package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-cz/textcase"
	"github.com/ivanvanderbyl/graphrag-go/pkg/extractors/entity"
	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := realMain(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func realMain(ctx context.Context) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: %s <document>", os.Args[0])
	}

	path := os.Args[1]
	documentText, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dbUser := ""
	dbPassword := ""
	dbUri := "bolt://localhost:7687" // scheme://host(:port) (default port is 7687)
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return err
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	} else {
		slog.Info("Driver connected to Memgraph")
	}

	indexes := []string{
		// "DROP GRAPH;",
		"CREATE INDEX ON :Policy(id);",
		// "CREATE INDEX ON :Policy(name);",
		"CREATE INDEX ON :Policy(embedding);",
		"CREATE INDEX ON :Person(id);",
		// "CREATE INDEX ON :Person(name);",
		"CREATE INDEX ON :Person(embedding);",
		"CREATE INDEX ON :Bill(id);",
		// "CREATE INDEX ON :Bill(name);",
		"CREATE INDEX ON :Bill(embedding);",
		"CREATE INDEX ON :Organization(id);",
		// "CREATE INDEX ON :Organization(name);",
		"CREATE INDEX ON :Organization(embedding);",
	}

	slog.Info("Extracting entities from document")
	openAILLM := llm.NewOpenAI(llm.WithCache(".cache"), llm.WithTemperature(0))

	extractor := entity.NewEntityExtractor(openAILLM)

	records, err := extractor.Extract(ctx, string(documentText))
	if err != nil {
		return err
	}

	// session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer session.Close(ctx)

	slog.Info("Creating Indexes")
	for _, index := range indexes {
		_, err = session.Run(ctx, index, nil)
		if err != nil {
			return err
		}
	}

	slog.Info("Creating Entities and Relationships")
	for _, record := range records {
		switch r := record.(type) {
		case *entity.Entity:
			slog.Info("Creating Entity", "name", r.Name)

			nodeType := pascalCase(r.Type())

			_, err := session.Run(ctx, "MERGE (e:$nodeType {id: $id, type: $type, name: $name, description: $description, embedding: $embedding}) RETURN e", map[string]interface{}{
				"id":          r.NodeID(),
				"nodeType":    nodeType,
				"type":        r.Type(),
				"name":        r.Name,
				"description": r.Description,
				"embedding":   r.Embedding,
			})
			if err != nil {
				return err
			}

		case *entity.Relationship:
			slog.Info("Creating Relationship", "from", r.Entity1, "relation", r.Keyword, "to", r.Entity2)

			fromNodeType := findRecordNodeTypes(records, r.Entity1)
			toNodeType := findRecordNodeTypes(records, r.Entity2)
			relationType := strings.ToUpper(r.Keyword)

			if relationType == "" {
				relationType = "RELATES"
			}

			if fromNodeType == "" || toNodeType == "" {
				// return fmt.Errorf("Could not find node types for relationship %s", r)

				slog.Warn("Could not find node types for relationship", "from", r.Entity1, "to", r.Entity2)
				continue
			}

			query := fmt.Sprintf("MATCH (from:$fromEntity {name: $from}), (to:$toEntity {name: $to}) MERGE (from)-[:%s {relation: $relation, weight: $weight}]->(to)", relationType)

			_, err = session.Run(ctx, query, map[string]interface{}{
				"from":       r.Entity1,
				"to":         r.Entity2,
				"relation":   r.Relation,
				"fromEntity": fromNodeType,
				"toEntity":   toNodeType,
				"weight":     r.Weight,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pascalCase(s string) string {
	return textcase.PascalCase(s)
}

func findRecordNodeTypes(records []entity.Record, relation string) string {
	for _, record := range records {
		switch r := record.(type) {
		case *entity.Entity:
			if r.Name == relation {
				return pascalCase(r.Type())
			}
		}
	}
	return ""
}
