package summarize_test

import (
	"context"
	"os"
	"testing"

	"github.com/ivanvanderbyl/graphrag-go/pkg/extractors/summarize"
	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/stretchr/testify/require"
)

func TestSummarizeExtractor_WithEntityNameKey(t *testing.T) {
	r := require.New(t)
	oai := llm.NewOpenAI(llm.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

	sum := summarize.NewSummarizeExtractor(oai)
	result, err := sum.Summarize(context.TODO(),
		[]string{"AUSTRALIAN"},
		[]string{
			"Australian refers to the people of Australia, indicating the geographical context of the statement.",
			"Australian refers to the people of Australia, who are the focus of efforts to reduce the cost of living.",
		},
	)
	r.NoError(err)
	r.NotNil(result)

	r.Equal("AUSTRALIAN", result.Items[0])
	t.Log(result.Description)
}

func TestSummarizeExtractor_WithSyntheticData(t *testing.T) {
	r := require.New(t)
	oai := llm.NewOpenAI(llm.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

	sum := summarize.NewSummarizeExtractor(oai, summarize.WithMaxSummaryLength(500))
	result, err := sum.Summarize(context.TODO(),
		[]string{"The Grand Canyon"},
		[]string{
			"The Grand Canyon is a steep-sided canyon carved by the Colorado River in Arizona, United States.",
			"It is contained within and managed by Grand Canyon National Park, the Hualapai Tribal Nation, and the Havasupai Tribe.",
			"The Grand Canyon is 277 miles long, up to 18 miles wide, and attains a depth of over a mile (6,093 feet or 1,857 meters).",
			"President Theodore Roosevelt was a major proponent of the preservation of the Grand Canyon area and visited it on numerous occasions to hunt and enjoy the scenery.",
			"The Grand Canyon is considered one of the Seven Natural Wonders of the World.",
			"The canyon offers picturesque views and is a popular destination for hikers, campers, and tourists.",
			"Geologically, the Grand Canyon reveals millions of years of Earth's history through its layered rock formations.",
		},
	)
	r.NoError(err)
	r.NotNil(result)

	r.Equal("The Grand Canyon", result.Items[0])
	t.Log(result.Description)
}
