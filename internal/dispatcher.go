package internal

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/diverged/uspt-go/internal/pipeline"
	"github.com/diverged/uspt-go/internal/utils"
	"github.com/diverged/uspt-go/types"
)

func Dispatcher(cfg *types.USPTGoConfig) (docChanOut <-chan *types.USPTGoDoc, errChanOut <-chan error, err error) {

	log := cfg.Logger

	log.Debug("Dispatcher called", "path", cfg.InputPath)

	docChan := make(chan *types.USPTGoDoc, 1000)
	errChan := make(chan error, 1000)

	zipFilePath := cfg.InputPath

	if filepath.Ext(zipFilePath) != ".zip" {
		err = errors.New("file is not a zip archive")
		errChan <- &types.USPTGoError{
			Err:     err,
			Skipped: true,
			Name:    filepath.Base(zipFilePath),
			Type:    "zip",
			Whence:  "file is not a zip archive",
		}
		close(errChan)
		return docChan, errChan, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Inspect the zip to determine the file format and schema version
		zipProfile, err := utils.InspectZip(zipFilePath)
		if err != nil {
			log.Error("zip file skipped due to error encountered while inspecting", "path", zipFilePath)
			err = errors.New("zip file skipped due to error encountered while inspecting")
			errChan <- &types.USPTGoError{
				Err:     err,
				Skipped: true,
				Name:    zipFilePath,
				Type:    "zip",
				Whence:  "while attempting to inspect the zip file",
			}
			close(errChan)
			close(docChan)
			// return nil, nil, errors.New("error encountered while inspecting zip file")
		}
		log.Debug("zip inspected successfully, proceeding with parsing logic", "path", zipFilePath)

		// Forward on based on zipProfile results
		switch zipProfile.OriginZip.ZipEntryExt {
		case ".xml":
			// Process XML files
			log.Debug("matched .xml zip entry extension", "path", zipProfile.OriginZip.ZipName)
			pipeline.XMLPipeline(zipProfile, cfg, docChan, errChan)

		case ".aps":
			// Process APS files
			log.Debug("matched .aps zip entry extension, handling for which is not yet implemented.", "path", zipProfile.OriginZip.ZipName)

		default:
			log.Error("Unknown file extension inside zip file", "path", zipFilePath, "extension", zipProfile.OriginZip.ZipEntryExt)
			errChan <- &types.USPTGoError{
				Err:     errors.New("unknown file extension mistakenly encountered within zip file"),
				Name:    zipProfile.OriginZip.ZipName,
				Type:    zipProfile.OriginZip.ZipEntryExt,
				Whence:  "while attempting to profile the zip",
				Skipped: true,
			}
			// Return exits the go routine on the unrecognized file extension, effectively skipping that zip file
			// return nil, nil, errors.New("unrecognized file extension found within zip")

		}
	}()
	wg.Wait()
	// Type-cast channels to receive-only for the return values
	docChanOut = docChan
	errChanOut = errChan

	return docChanOut, errChanOut, nil
}
