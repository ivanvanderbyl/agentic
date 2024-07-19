package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDocument = "(\"entity\"<|>\"Michelle Ananda-Rajah\"<|>\"person\"<|>\"Michelle Ananda-Rajah is a Member of Parliament (MP) from the ALP, representing the electorate of Higgins. She delivered a speech addressing environmental issues and legislation.\"##  \n(\"entity\"<|>\"The ALP\"<|>\"organization\"<|>\"The ALP (Australian Labor Party) is a political party in Australia representing progressive policies. Michelle Ananda-Rajah is a member of this party.\"##  \n(\"entity\"<|>\"Higgins\"<|>\"electorate\"<|>\"Higgins is an electorate in Australia represented by Michelle Ananda-Rajah, MP.\"##  \n(\"entity\"<|>\"The Liberal Government\"<|>\"organization\"<|>\"The Liberal Government refers to the previous governing party in Australia, criticized for suppressing the state of the environment report and neglecting environmental protection.\"##  \n(\"entity\"<|>\"Climate Change\"<|>\"policy\"<|>\"Climate Change refers to the environmental policy issue related to the alteration of the planet's climate patterns due to human activities, emphasized as a pressing concern in the speech.\"##  \n(\"entity\"<|>\"EPBC Act\"<|>\"policy\"<|>\"The EPBC (Environment Protection and Biodiversity Conservation) Act is a legislative framework for environmental protection in Australia, criticized for being ineffective.\"##  \n(\"entity\"<|>\"Nature Repair Market Bill\"<|>\"bill\"<|>\"The Nature Repair Market Bill is a legislative measure passed to expand environmental protections, including aspects such as the water trigger for unconventional gas.\"##  \n(\"entity\"<|>\"Nature Positive (Environment Protection Australia) Bill 2024\"<|>\"bill\"<|>\"The Nature Positive Bill 2024 is a legislative package aiming to enhance environmental protection and improve information transparency.\"##  \n(\"entity\"<|>\"Environment Protection Authority (EPA)\"<|>\"organization\"<|>\"The EPA (Environment Protection Authority) is an independent statutory body dedicated to enforcing and regulating environmental standards across Australia.\"##  \n(\"entity\"<|>\"Environment Information Australia (EIA)\"<|>\"organization\"<|>\"The EIA (Environment Information Australia) is an organization tasked with gathering and providing environmental data to support decision-making.\"##  \n(\"entity\"<|>\"Minister Plibersek\"<|>\"person\"<|>\"Minister Plibersek is a government minister commended for driving the environmental legislation discussed in the speech.\"##  \n(\"entity\"<|>\"Professor Samuel\"<|>\"person\"<|>\"Professor Graham Samuel is a critic of the EPBC Act and provided a review calling for stricter environmental protection measures.\"##  \n(\"entity\"<|>\"State of the Environment Report\"<|>\"policy\"<|>\"The State of the Environment Report is an official assessment of Australia's environmental condition. The 2021 report was suppressed by the previous government due to its damning findings.\"##  \n(\"entity\"<|>\"Environmental Defenders Office\"<|>\"organization\"<|>\"The Environmental Defenders Office is a legal service focused on environmental justice, which the current opposition leader intends to defund.\"##  \n(\"entity\"<|>\"WWF (World Wildlife Fund)\"<|>\"organization\"<|>\"The WWF (World Wildlife Fund) is a global environmental organization that supports the establishment of the EPA as a potential game changer.\"##  \n(\"entity\"<|>\"The Albanese Government\"<|>\"organization\"<|>\"The Albanese Government is the current governing party in Australia, driving the nature-positive plan for a sustainable future.\"##  \n(\"relationship\"<|>\"Michelle Ananda-Rajah\"<|>\"The ALP\"<|>\"Michelle Ananda-Rajah is a member of the ALP and represents its policies and goals.\"<|>10<|>MEMBER_OF)##  \n(\"relationship\"<|>\"Michelle Ananda-Rajah\"<|>\"Higgins\"<|>\"Michelle Ananda-Rajah represents the electorate of Higgins in the Australian Parliament.\"<|>10<|>REPRESENTS)##  \n(\"relationship\"<|>\"The Liberal Government\"<|>\"State of the Environment Report\"<|>\"The previous Liberal Government suppressed the State of the Environment Report due to its critical findings.\"<|>9<|>SUPPRESSED)##  \n(\"relationship\"<|>\"EPBC Act\"<|>\"Professor Samuel\"<|>\"Professor Samuel criticized the EPBC Act for being ineffective and called it 'gobbledygook.'\"<|>8<|>CRITICIZED)##  \n(\"relationship\"<|>\"Nature Repair Market Bill\"<|>\"EPBC Act\"<|>\"The Nature Repair Market Bill was introduced to address the deficiencies found in the EPBC Act.\"<|>9<|>REFORMS)##  \n(\"relationship\"<|>\"Nature Positive (Environment Protection Australia) Bill 2024\"<|>\"EPA\"<|>\"The Nature Positive Bill 2024 aims to establish the Environment Protection Authority (EPA) as a regulator.\"<|>9<|>ESTABLISHES)##  \n(\"relationship\"<|>\"EPA\"<|>\"EIA\"<|>\"The EPA will work with the EIA to better use environmental data for planning and decision-making.\"<|>8<|>COLLABORATES)##  \n(\"relationship\"<|>\"Minister Plibersek\"<|>\"The Albanese Government\"<|>\"Minister Plibersek is part of the Albanese Government and is driving the environmental legislation.\"<|>9<|>PART_OF)##  \n(\"relationship\"<|>\"The Albanese Government\"<|>\"EPA\"<|>\"The Albanese Government is establishing the EPA to enhance environmental protections.\"<|>10<|>ESTABLISHES)##  \n(\"relationship\"<|>\"The Albanese Government\"<|>\"Nature Positive (Environment Protection Australia) Bill 2024\"<|>\"The Albanese Government is promoting the Nature Positive Bill 2024 to improve environmental protections.\"<|>9<|>PROMOTES)##  \n(\"relationship\"<|>\"Environmental Defenders Office\"<|>\"State of the Environment Report\"<|>\"The opposition leader's plan to defund the Environmental Defenders Office contrasts with the need for transparency shown in the State of the Environment Report.\"<|>7<|>CONTRASTS)##  \n(\"relationship\"<|>\"WWF\"<|>\"EPA\"<|>\"The WWF supports the establishment of the EPA as a significant step in environmental protection.\"<|>8<|>SUPPORTS)<|COMPLETE|>"

func TestEntityParsing(t *testing.T) {
	a := assert.New(t)

	extractor := NewEntityExtractor(nil)
	records := extractor.processResults(testDocument)

	if len(records) != 28 {
		t.Fatalf("Expected 28 records, got %d", len(records))
	}

	person := records[0].(*Entity)
	a.Equal("Michelle Ananda-Rajah", person.Name)

	var relation Relationship
	for _, record := range records {
		if r, ok := record.(*Relationship); !ok {
			continue
		} else {
			relation = *r
			break
		}
	}

	a.Equal("Michelle Ananda-Rajah", relation.Entity1)
	a.Equal("The ALP", relation.Entity2)
	a.Equal("Michelle Ananda-Rajah is a member of the ALP and represents its policies and goals.", relation.Relation)
	a.Equal(10, relation.Weight)
	a.Equal("MEMBER_OF", relation.Keyword)
}

func findRecordNodeTypes(records []Record, relation string) string {
	for _, record := range records {
		switch r := record.(type) {
		case *Entity:
			if r.Name == relation {
				return r.Type()
			}
		}
	}
	return ""
}
