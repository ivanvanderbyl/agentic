package model

type (
	Document struct {
		Identified
		Text string `json:"text"`

		// TextUnits is a list of TextUnit IDs that are part of this document.
		TextUnits []*TextUnit `json:"text_units,omitempty"`

		// ExtractedEntities is a list of entities that were extracted from the document.
		ExtractedEntities ExtractedEntities `json:"extracted_entities,omitempty"`
	}

	ExtractEntity     string
	ExtractedEntities []ExtractEntity
)
