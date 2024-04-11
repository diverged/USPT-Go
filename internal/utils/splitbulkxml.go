package utils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/diverged/uspt-go/types"
)

// BulkXMLSplitter processes a zip file containing bulk XML documents, sending individual documents to a channel.
func BulkXMLSplitter(bulkZip *types.USPTGoMetadata, splitXMLDocChan chan<- *types.USPTGoDoc, errChan chan<- error, log types.Logger) {

	zipInfo := bulkZip.OriginZip

	log.Info("BulkXMLSplitter starting for zip file", "Zip File", zipInfo.ZipName)

	log.Info("Calling zip.OpenReader", "Zip File", zipInfo.ZipName)
	zipReader, err := zip.OpenReader(zipInfo.ZipPath)
	if err != nil {
		log.Error("Failed to open zip file for processing", "Zip File Path", zipInfo.ZipPath, "Error", err)
		errChan <- &types.USPTGoError{
			Skipped: true,
			Name:    zipInfo.ZipName,
			Type:    "zip",
			Whence:  "opening zip file",
			Err:     err,
		}
		return
	}
	defer zipReader.Close()

	for _, zipEntry := range zipReader.File {

		if zipEntry.FileInfo().IsDir() {
			continue // Skip directories
		}

		log.Info("Starting splitting process on bulk file entry inside the zip", "Zip entry within", zipInfo.ZipName, "with file extension", zipInfo.ZipEntryExt)
		f, err := zipEntry.Open()
		if err != nil {
			errChan <- &types.USPTGoError{
				Skipped: true,
				Name:    zipInfo.ZipName,
				Type:    "zip entry",
				Whence:  "attempting to open zip entry for reading",
				Err:     err,
			}
			continue
		}

		// * Now can call processXMLDocument on the Zip Entry
		processXMLDocument(bulkZip, zipEntry, f, splitXMLDocChan, errChan, log)

		f.Close() // Close the zip entry file handle after processing

		log.Info("splitting of bulk file completed, entry is now closed", "Zip Name", zipInfo.ZipName, "timestamp", time.Now())
	}
}

func processXMLDocument(zipInfo *types.USPTGoMetadata, zipEntry *zip.File, f io.ReadCloser, splitXMLDocChan chan<- *types.USPTGoDoc, errChan chan<- error, log types.Logger) {
	bufferedReader := bufio.NewReader(f)
	var buffer bytes.Buffer
	var inXMLDocument bool
	xmlStartTag := []byte("<?xml")
	documentIndex := 0 // Initialize a counter for each XML document

	for {
		line, err := bufferedReader.ReadBytes('\n')
		if bytes.Contains(line, xmlStartTag) && inXMLDocument {
			// End of the current XML document has been reached.  Trim the trailing newline character from the buffer before sending the document.
			trimmedBuffer := bytes.TrimSpace(buffer.Bytes())
			sendDocument(zipInfo, zipEntry, documentIndex, bytes.NewBuffer(trimmedBuffer), splitXMLDocChan, log)
			buffer.Reset()
			inXMLDocument = false
			documentIndex++ // Increment the counter for the next document
		}
		if bytes.Contains(line, xmlStartTag) {
			inXMLDocument = true // Start of a new XML document
		}
		if inXMLDocument {
			buffer.Write(line)
		}
		if err == io.EOF {
			if inXMLDocument {
				// End of the last XML document in the file.  Trim the trailing newline character from the buffer before sending the document.
				trimmedBuffer := bytes.TrimSpace(buffer.Bytes())
				sendDocument(zipInfo, zipEntry, documentIndex, bytes.NewBuffer(trimmedBuffer), splitXMLDocChan, log)
			}
			break
		} else if err != nil {
			// log.Error("Error encountered while reading zip entry", zap.String("zipEntryName", zipEntry.Name), zap.Error(err))
			// errChan <- err
			errChan <- &types.USPTGoError{
				Skipped: true,
				Name:    zipEntry.Name,
				Type:    "bulk xml zip",
				Whence:  "attempting to read zip entry",
				Err:     err,
			}
			break
		}
	}
}

func sendDocument(zipInfo *types.USPTGoMetadata, zipEntry *zip.File, documentIndex int, buffer *bytes.Buffer,
	splitXMLDocChan chan<- *types.USPTGoDoc, log types.Logger) {

	filename := fmt.Sprintf("%s-%d.xml", strings.TrimSuffix(zipEntry.Name, filepath.Ext(zipEntry.Name)), documentIndex)

	zipInfo.OriginZip.IndexName = filename

	// Simple indicator of []byte slice integrity on way out
	if buffer.Bytes()[len(buffer.Bytes())-1] != '>' {
		log.Error("Document does not end with a closing tag", "filename", filename)
	}

	// Create a deep copy of the buffer to send through the channel
	originalXML := buffer.Bytes()
	copiedXML := make([]byte, len(originalXML)) // instantiate the new variable with a byte slice sized to match the source byte slice being copied
	copy(copiedXML, originalXML)                // Send 'copiedXML' through the channel

	// Send the document to the splitXMLDocChan
	switch zipInfo.DocumentType {
	case "trademark":
		splitXMLDocChan <- &types.USPTGoDoc{
			USPTGoMetadata: *zipInfo,
			Trademark: types.Trademark{
				RawSplitDoc: copiedXML,
			},
		}
	default:
		splitXMLDocChan <- &types.USPTGoDoc{
			USPTGoMetadata: *zipInfo,
			RawSplitDoc:    copiedXML,
		}
	}

	log.Debug("BulkXMLSplitter: doc => splitXMLDocChan", "filename", filename)
}
