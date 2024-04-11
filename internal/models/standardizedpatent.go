package models

/*
type StandardizedPatentContents struct {
	PatentMetadata StdPatentMetadata
	BiblioData     StdBiblioData
}

type StdPatentMetadata struct {
	// "Meta-" fields correspond to XML Namespace attributes
	MetaLang         string `json:"meta_lang"`
	MetaFileName     string `json:"meta_file_name"`
	MetaFileStatus   string `json:"meta_file_status"`
	MetaFileType     string `json:"meta_file_type"`
	MetaCountry      string `json:"meta_country"`
	MetaDateProduced string `json:"meta_date_produced"`
	MetaDatePubl     string `json:"meta_date_publ"`
	// NumberOfClaims and InventionTitle also treated as a "Meta-" as they don't intuitively fit with the other Biblio fields
	MetaNumberOfClaims int    `json:"number_of_claims"`
	MetaInventionTitle string `json:"invention_title"`
}

// Interface for standardizing both Grant and Application versions of Bibliographic Data Fields
type StdPatentContentConverter interface {
	ToStdPatentContents() StandardizedPatentContents
}

type StdBiblioData struct {
	UsBibliographicData StdUsBibliographicData `json:"us_bibliographic_data"`
}

// ToStdPatentContents converts Grant data to StdBiblioData
func (bg XMLGrantBibliographicData) ToStdPatentContents() StandardizedPatentContents {
	// Extracting metadata
	metadata := StdPatentMetadata{
		MetaLang:           bg.MetaLang,
		MetaFileName:       bg.MetaFileName,
		MetaFileStatus:     bg.MetaStatus,
		MetaFileType:       bg.MetaFileType,
		MetaCountry:        bg.MetaCountry,
		MetaDateProduced:   bg.MetaDateProduced,
		MetaDatePubl:       bg.MetaDatePubl,
		MetaNumberOfClaims: bg.UsBibliographicData.NumberOfClaims,
		MetaInventionTitle: bg.UsBibliographicData.InventionTitle.Text,
	}

	// Extracting bibliographic data
	biblioData := StdBiblioData{
		UsBibliographicData: StdUsBibliographicData{
			PublicationReference: StdPublicationReference{
				DocumentID: StdDocumentID{
					Country:   bg.UsBibliographicData.PublicationReference.DocumentID.Country,
					DocNumber: bg.UsBibliographicData.PublicationReference.DocumentID.DocNumber,
					KindCode:  bg.UsBibliographicData.PublicationReference.DocumentID.KindCode,
					Date:      bg.UsBibliographicData.PublicationReference.DocumentID.Date,
				},
			},
			ApplicationReference: StdApplicationReference{
				ApplType: bg.UsBibliographicData.ApplicationReference.ApplType,
				DocumentID: StdDocumentID{
					Country:   bg.UsBibliographicData.ApplicationReference.DocumentID.Country,
					DocNumber: bg.UsBibliographicData.ApplicationReference.DocumentID.DocNumber,
					Date:      bg.UsBibliographicData.ApplicationReference.DocumentID.Date,
				},
			},
			ClassificationNational: StdClassificationNational{
				Country:               bg.UsBibliographicData.ClassificationNational.Country,
				MainClassification:    bg.UsBibliographicData.ClassificationNational.MainClassification,
				FurtherClassification: bg.UsBibliographicData.ClassificationNational.FurtherClassification,
			},
		},
	}

	// Constructing the final StandardizedPatentContents struct
	return StandardizedPatentContents{
		PatentMetadata: metadata,
		BiblioData:     biblioData,
	}
}

// ToStdPatentContents converts Grant data to StdBiblioData
func (ba XMLApplicationBibliographicData) ToStdPatentContents() StandardizedPatentContents {
	// Extracting metadata
	metadata := StdPatentMetadata{
		MetaLang:           ba.MetaLang,
		MetaFileName:       ba.MetaFileName,
		MetaFileStatus:     ba.MetaStatus,
		MetaFileType:       ba.MetaFileType,
		MetaCountry:        ba.MetaCountry,
		MetaDateProduced:   ba.MetaDateProduced,
		MetaDatePubl:       ba.MetaDatePubl,
		MetaNumberOfClaims: ba.UsBibliographicData.NumberOfClaims,
		MetaInventionTitle: ba.UsBibliographicData.InventionTitle.Text,
	}

	// Extracting bibliographic data
	biblioData := StdBiblioData{
		UsBibliographicData: StdUsBibliographicData{
			PublicationReference: StdPublicationReference{
				DocumentID: StdDocumentID{
					Country:   ba.UsBibliographicData.PublicationReference.DocumentID.Country,
					DocNumber: ba.UsBibliographicData.PublicationReference.DocumentID.DocNumber,
					KindCode:  ba.UsBibliographicData.PublicationReference.DocumentID.KindCode,
					Date:      ba.UsBibliographicData.PublicationReference.DocumentID.Date,
				},
			},
			ApplicationReference: StdApplicationReference{
				ApplType: ba.UsBibliographicData.ApplicationReference.ApplType,
				DocumentID: StdDocumentID{
					Country:   ba.UsBibliographicData.ApplicationReference.DocumentID.Country,
					DocNumber: ba.UsBibliographicData.ApplicationReference.DocumentID.DocNumber,
					Date:      ba.UsBibliographicData.ApplicationReference.DocumentID.Date,
				},
			},
			ClassificationNational: StdClassificationNational{
				Country:               ba.UsBibliographicData.ClassificationNational.Country,
				MainClassification:    ba.UsBibliographicData.ClassificationNational.MainClassification,
				FurtherClassification: ba.UsBibliographicData.ClassificationNational.FurtherClassification,
			},
		},
	}

	// Constructing the final StandardizedPatentContents struct
	return StandardizedPatentContents{
		PatentMetadata: metadata,
		BiblioData:     biblioData,
	}
}

// Breakdown of StdBiblioData to avoid nested anonymous structs
type StdUsBibliographicData struct {
	PublicationReference   StdPublicationReference   `json:"publication_reference"`
	ApplicationReference   StdApplicationReference   `json:"application_reference"`
	ClassificationNational StdClassificationNational `json:"classification_national"`
}

type StdPublicationReference struct {
	DocumentID StdDocumentID `json:"document_id"`
}

type StdApplicationReference struct {
	ApplType   string        `json:"appl_type"`
	DocumentID StdDocumentID `json:"document_id"`
}

type StdClassificationNational struct {
	Country               string `json:"country"`
	MainClassification    string `json:"main_classification"`
	FurtherClassification string `json:"further_classification"`
}

type StdDocumentID struct {
	Country   string `json:"country"`
	DocNumber string `json:"doc_number"`
	KindCode  string `json:"kind_code"`
	Date      string `json:"date"`
}
*/
