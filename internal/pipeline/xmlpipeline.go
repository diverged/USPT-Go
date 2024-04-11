package pipeline

import (
	"github.com/diverged/uspt-go/internal/parsers/xmlparser"
	"github.com/diverged/uspt-go/internal/transformtext"
	"github.com/diverged/uspt-go/internal/utils"
	"github.com/diverged/uspt-go/types"
)

// XMLPipeline is the processing logic flow for bulk XML patent files of both Grant and Application types.
func XMLPipeline(zipProfile *types.USPTGoMetadata, cfg *types.USPTGoConfig, docChan chan<- *types.USPTGoDoc, errChan chan<- error) {

	bulkZip := zipProfile.OriginZip
	log := cfg.Logger
	log.Info("Starting XMLPipeline", "Bulk Zip File", bulkZip.ZipName)

	// Create blocking channels
	splitXMLDocChan := make(chan *types.USPTGoDoc, 100)  // BulkXMLSplitter() => splitXMLDocChan => XMLParser()
	parsedXMLDocChan := make(chan *types.USPTGoDoc, 100) // XMLParser() => parsedXMLDocChan
	transDocChan := make(chan *types.USPTGoDoc, 100)     // XMLParser() => transDocChan

	// Start BulkXMLSplitter() in goroutine
	go func() {
		log.Info("Initializing XML Bulk Splitter", "Splitting", bulkZip.ZipName)
		defer close(splitXMLDocChan)
		utils.BulkXMLSplitter(zipProfile, splitXMLDocChan, errChan, log)
	}()

	// Start ParseXMLPatent() in go routine
	go func() {
		defer close(parsedXMLDocChan)

		switch zipProfile.DocumentType {
		case "application", "grant":
			log.Info("Initializing parsing of split XML Patent docs")
			xmlparser.ParseXMLPatent(cfg, splitXMLDocChan, parsedXMLDocChan, errChan, log)
		// case "trademark":
		// log.Info("Initializing parsing of split XML trademark docs")
		default:
			log.Error("XMLPipeline() couldn't match DocumentType when assigning a parser")
		}
	}()

	// Start TranslatePatentXmlToHtml() in go routine to translate XML to HTML
	go func() {
		log.Info("Initializing XML to HTML translation")
		defer close(transDocChan)
		transformtext.TranslatePatentXmlToHtml(parsedXMLDocChan, transDocChan, errChan, log)
	}()

	// Forward on the channel contents
	go func() {

		defer close(docChan)
		defer close(errChan)
		for doc := range transDocChan {
			docChan <- doc
		}
	}()

}
