package models

import (
	"encoding/xml"
)

type XMLClaim struct {
	XMLName   xml.Name `xml:"claim"`
	ID        string   `xml:"id,attr"`
	Num       string   `xml:"num,attr"`
	ClaimText XMLClaimText
}

type XMLClaimText struct {
	XMLName  xml.Name              `xml:"claim-text"`
	Text     string                `xml:",innerxml"`
	Elements []XMLClaimTextElement `xml:"claim-text"`
	ClaimRef XMLClaimRef           `xml:"claim-ref"`
}

type XMLClaimRef struct {
	XMLName xml.Name `xml:"claim-ref"`
	IDRef   string   `xml:"idref,attr"`
}

type XMLClaimTextElement struct {
	XMLName   xml.Name `xml:"claim-text"`
	Text      string   `xml:",innerxml"`
	ClaimRefs []struct {
		XMLName xml.Name `xml:"claim-ref"`
		IDRef   string   `xml:"idref,attr"`
	} `xml:"claim-ref"`
}

type Claim struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Text      []string  `json:"text"`
	ClaimTree ClaimTree `json:"claimTree"`
	ChildIds  []string  `json:"childIds,omitempty"`
}

type ClaimTree struct {
	ParentIds      []string `json:"parentIds,omitempty"`
	ParentCount    int      `json:"parentCount"`
	ChildCount     int      `json:"childCount"`
	ClaimTreeLevel int      `json:"claimTreelevel"`
}
