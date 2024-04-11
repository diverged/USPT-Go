package types

import (
	"encoding/xml"

	"github.com/diverged/uspt-go/internal/models"
)

type Patent struct {
	XMLName             xml.Name            `xml:"-" json:"-"` // `xml:"us-patent-grant"` OR `xml:"us-patent-application"`
	MetaLang            string              `xml:"lang,attr" json:"lang"`
	MetaDtdVersion      string              `xml:"dtd-version,attr" json:"dtd-version"`
	MetaFileName        string              `xml:"file,attr" json:"file-name"`
	MetaStatus          string              `xml:"status,attr" json:"status"`
	MetaFileType        string              `xml:"id,attr" json:"id"`
	MetaCountry         string              `xml:"country,attr" json:"country"`
	MetaDateProduced    string              `xml:"date-produced,attr" json:"date-produced"`
	MetaDatePubl        string              `xml:"date-publ,attr" json:"date-publ"`
	UsBibliographicData UsBibliographicData `xml:"-" json:"-"` // `xml:"us-bibliographic-data-grant"` OR `xml:"us-bibliographic-data-application"`
	Description         struct {
		Content string `xml:",innerxml"`
	} `xml:"description"`
	Abstract struct {
		Content string `xml:",innerxml"`
	} `xml:"abstract"`
	Claims struct {
		Content string `xml:",innerxml"`
	} `xml:"claims"`
	StructuredClaims []*models.Claim
}

type UsBibliographicData struct {
	XMLName              xml.Name `xml:"-" json:"-"` // `xml:"us-bibliographic-data-grant"` OR `xml:"us-bibliographic-data-application"`
	PublicationReference struct {
		DocumentID struct {
			Country   string `xml:"country"`
			DocNumber string `xml:"doc-number"`
			KindCode  string `xml:"kind"`
			Date      string `xml:"date"`
		} `xml:"document-id"`
	} `xml:"publication-reference"`
	ApplicationReference struct {
		ApplType   string `xml:"appl-type,attr"`
		DocumentID struct {
			Country   string `xml:"country"`
			DocNumber string `xml:"doc-number"`
			Date      string `xml:"date"`
		} `xml:"document-id"`
	} `xml:"application-reference"`
	ClassificationNational struct {
		Country               string `xml:"country"`
		MainClassification    string `xml:"main-classification"`
		FurtherClassification string `xml:"further-classification"`
	} `xml:"classification-national"`
	InventionTitle struct {
		Content string `xml:",innerxml"`
		Text    string `xml:",chardata"`
		ID      string `xml:"id,attr"`
	} `xml:"invention-title"`
	NumberOfClaims int `xml:"number-of-claims"`
}
