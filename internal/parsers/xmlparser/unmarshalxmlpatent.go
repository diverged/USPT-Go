package xmlparser

import (
	"encoding/xml"
	"errors"

	"github.com/diverged/uspt-go/types"
)

func UnmarshalXmlPatent(rawSplitDoc []byte, patentDocType string, errChan chan<- error, log types.Logger) (types.Patent, error) {
	var patent types.Patent

	log.Debug("UnmarshalXmlPatent has been called")

	switch patentDocType {
	case "application":
		patent.XMLName = xml.Name{Local: "us-patent-application"}
		patent.UsBibliographicData.XMLName = xml.Name{Local: "us-bibliographic-data-application"}
	case "grant":
		patent.XMLName = xml.Name{Local: "us-patent-grant"}
		patent.UsBibliographicData.XMLName = xml.Name{Local: "us-bibliographic-data-grant"}
	default:
		log.Error("Unknown document type: %s", patentDocType)
		err := errors.New("unknown document type when attempting to unmarshal xml patent")
		return types.Patent{}, err
	}

	if err := xml.Unmarshal(rawSplitDoc, &patent); err != nil {
		log.Error("Error unmarshaling xml patent", "error", err)
		return types.Patent{}, err
	}

	return patent, nil
}
