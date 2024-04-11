package xmlparser

import (
	// "bytes"

	"errors"
	"fmt"

	// "io"
	"strings"

	"github.com/diverged/uspt-go/types"
)

func ParseXMLPatent(cfg *types.USPTGoConfig, splitXMLDocChan <-chan *types.USPTGoDoc, parsedXMLDocChan chan<- *types.USPTGoDoc, errChan chan<- error, log types.Logger) {

	log.Debug("ParseXMLPatent has been invoked")

	for doc := range splitXMLDocChan {

		// fmt.Println(string(doc.RawSplitDoc))

		// log.Debug("Parsing split XML document", zap.String("DocName", doc.GoUSPTOMetadata.SplitterIndex["DocIndexOfZip"]))
		var (
			happyParser = true            // Abort flag to enable skipping the document
			parseErrors []error           // Accumulates any encountered parsing errors
			rawSplitDoc = doc.RawSplitDoc // Extract the raw XML document from the XMLDocument
		)

		// * Initial Unmarshaling

		unmarshaledPatent, err := UnmarshalXmlPatent(rawSplitDoc, doc.USPTGoMetadata.DocumentType, errChan, log)
		if err != nil {
			errChan <- &types.USPTGoError{
				Err:     err,
				Name:    doc.USPTGoMetadata.OriginZip.IndexName,
				Type:    "xml patent",
				Whence:  "unmarshaling XML document",
				Skipped: true,
			}
			continue
		}

		doc.Patent = unmarshaledPatent

		// fmt.Println(doc.Patent.Description.Content)

		// * Map the Claims Tree
		structuredClaims, err := ParseStructuredClaims(rawSplitDoc, log)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to parse structured claims from extracted xml claims []byte slice: %w", err))
			happyParser = false
		}

		// fmt.Println(structuredClaims)

		// Assign the structured claims to the doc
		doc.Patent.StructuredClaims = structuredClaims

		// * If parser is not happy, collect the parsing error(s) and report the skipped document to errChan
		if !happyParser {
			combinedError := combineErrors(parseErrors)
			errChan <- &types.USPTGoError{
				Err:     combinedError,
				Name:    doc.USPTGoMetadata.OriginZip.IndexName,
				Type:    "xml patent",
				Whence:  "parsing XML document",
				Skipped: true,
			}
			// `continue` bypasses the remaining code in the loop and starts the next iteration, effectively blocking the document from ever being sent into parsedXMLDocChan
			continue
		}

		// * Send the parsed XML document to the channel
		parsedXMLDocChan <- doc
		log.Debug("ParseXMLElements: doc => parsedXMLDocChan", "DocName", doc.USPTGoMetadata.OriginZip.IndexName)

	} // End of range over channel
}

// Helper to combine the slice of errors into a string
func combineErrors(parseErrors []error) error {
	var errMsgs []string
	for _, err := range parseErrors {
		errMsgs = append(errMsgs, err.Error())
	}
	return errors.New(strings.Join(errMsgs, "; "))
}

/* func ParseXMLPatent(cfg *types.USPTGoConfig, splitXMLDocChan <-chan *types.USPTGoDoc, parsedXMLDocChan chan<- *types.USPTGoDoc, errChan chan<- error, log types.Logger) {

	log.Debug("ParseXMLPatent has been invoked")

	for doc := range splitXMLDocChan {

		// fmt.Println(string(doc.RawSplitDoc))

		// log.Debug("Parsing split XML document", zap.String("DocName", doc.GoUSPTOMetadata.SplitterIndex["DocIndexOfZip"]))
		var (
			happyParser = true            // Abort flag to enable skipping the document
			parseErrors []error           // Accumulates any encountered parsing errors
			rawSplitDoc = doc.RawSplitDoc // Extract the raw XML document from the XMLDocument
		)

		// * Extract the Abstract, Description and Claims
		// extractedTextBytes, err := ExtractTextBytes(rawSplitDoc)
		// if err != nil {
		// 	parseErrors = append(parseErrors, fmt.Errorf("failed attempting to extract byte slices from XML: %w", err))
		// 	happyParser = false
		// }

		// TEMP DEBUG


		type TestInnerXml struct {
			Description struct {
				Raw []byte `xml:",innerxml"`
			} `xml:"description"`
			Abstract struct {
				Raw []byte `xml:",innerxml"`
			} `xml:"abstract"`
			Claims struct {
				Raw []byte `xml:",innerxml"`
			} `xml:"claims"`
		}

		var extractedTextBytes TestInnerXml
		err := xml.Unmarshal(rawSplitDoc, &extractedTextBytes)
		if err != nil {
			log.Error("Error unmarshalling InnerXml:", err)
		}

		fmt.Println(string(extractedTextBytes.Description.Raw))
		fmt.Println(string(extractedTextBytes.Abstract.Raw))
		// fmt.Println(string(extractedTextBytes.Claims.Raw))

		// TEMP DEBUG

		// * Map the Claims Tree
		// ! Needs the <claims> tags to be present in the XML
		//structuredClaims, err := ParseStructuredClaims(extractedTextBytes.Claims, log) //
		structuredClaims, err := ParseStructuredClaims(rawSplitDoc, log) //
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to parse structured claims from extracted xml claims []byte slice: %w", err))
			happyParser = false
		}

		// fmt.Println(structuredClaims)

		// Assign the structured claims to the doc
		doc.Patent.Claims.StructuredClaims = structuredClaims

		// * Format text fields per cfg

		var formatter TextFormatter
		switch cfg.OutputFormats.MainTextFields {
		case "html":
			formatter = HTMLFormatter{}
		case "markdown":
			formatter = MarkdownFormatter{}
		default:
			formatter = PlaintextFormatter{}
		}

		formattedDescription, err := formatter.FormatText(extractedTextBytes.Description.Raw)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to format description text: %w", err))
			happyParser = false
		}
		//fmt.Println(formattedDescription)

		formattedAbstract, err := formatter.FormatText(extractedTextBytes.Abstract.Raw)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to format abstract text: %w", err))
			happyParser = false
		}

		formattedClaimsText, err := formatter.FormatText(extractedTextBytes.Claims.Raw)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to format claims as text: %w", err))
			happyParser = false
		}

		// Assign the formatted text to the doc
		doc.Patent.Description = formattedDescription
		doc.Patent.Abstract = formattedAbstract
		doc.Patent.Claims.ClaimsText = formattedClaimsText

		// * Parse Bilbiographic data and XML Namespace metadata

		// stdContents is a struct which standardizes differences in bibliographic data between application and grant XML
		var stdContents models.StandardizedPatentContents

		switch doc.USPTGoMetadata.DocumentType {
		case "application":
			// Create a new decoder to parse an Application's Bibliographic data
			appData := &models.XMLApplicationBibliographicData{}
			if err := xml.Unmarshal(rawSplitDoc, appData); err != nil {

				parseErrors = append(parseErrors, fmt.Errorf("failed to unmarshal XML into XMLApplicationBibliographicData: %w", err))
				happyParser = false

			}
			// Standardize formatting of the Namespace Attribute Metadata and Bibliographic Data
			stdContents = appData.ToStdPatentContents()
		case "grant":
			// Create a new decoder to parse a Grant's Bibliographic data
			grantData := &models.XMLGrantBibliographicData{}
			if err := xml.Unmarshal(rawSplitDoc, grantData); err != nil {

				parseErrors = append(parseErrors, fmt.Errorf("failed to unmarshal XML into XMLGrantBibliographicData: %w", err))
				happyParser = false

			}
			// Standardize formatting of the Namespace Attribute Metadata and Bibliographic Data
			stdContents = grantData.ToStdPatentContents()
		default:

			parseErrors = append(parseErrors, fmt.Errorf("failed to match DocumentType to either \"application\" or \"grant\" during XML Parsing"))
			happyParser = false

		}

		// Assign standardized PatentMetadata and Bibliographic Data to doc
		doc.Patent.MetaData = stdContents.PatentMetadata
		doc.Patent.BiblioData = stdContents.BiblioData

		// * If parser is not happy, collect the parsing error(s) and report the skipped document to errChan
		if !happyParser {
			combinedError := combineErrors(parseErrors)
			errChan <- &types.USPTGoError{
				Err:     combinedError,
				Name:    doc.USPTGoMetadata.OriginZip.IndexName,
				Type:    "xml patent",
				Whence:  "parsing XML document",
				Skipped: true,
			}
			// `continue` bypasses the remaining code in the loop and starts the next iteration, effectively blocking the document from ever being sent into parsedXMLDocChan
			continue
		}

		// * Send the parsed XML document to the channel
		parsedXMLDocChan <- doc
		log.Debug("ParseXMLElements: doc => parsedXMLDocChan", "DocName", doc.USPTGoMetadata.OriginZip.IndexName)

	} // End of range over channel
}

// Helper to combine the slice of errors into a string
func combineErrors(parseErrors []error) error {
	var errMsgs []string
	for _, err := range parseErrors {
		errMsgs = append(errMsgs, err.Error())
	}
	return errors.New(strings.Join(errMsgs, "; "))
}
*/
