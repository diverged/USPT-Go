package transformtext

import (
	"github.com/diverged/uspt-go/types"
)

func TranslatePatentXmlToHtml(parsedXmlDocChan <-chan *types.USPTGoDoc, transDocChan chan<- *types.USPTGoDoc, errChan chan<- error, log types.Logger) {
	for doc := range parsedXmlDocChan {
		// Translate the inner XML content to HTML
		htmlDescription, err := InnerXmlToHtml([]byte(doc.Patent.Description.Content))
		if err != nil {
			errChan <- &types.USPTGoError{
				Skipped: true,
				Name:    doc.USPTGoMetadata.OriginZip.IndexName,
				Type:    "xml translation",
			}
		}
		doc.Patent.Description.Content = htmlDescription

		transDocChan <- doc
	}
}
