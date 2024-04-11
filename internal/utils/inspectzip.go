package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/diverged/uspt-go/types"
)

// TODO - Just return errors created here using constructor, then within Controller, drop them into errChan

func InspectZip(zipFilePath string) (*types.USPTGoMetadata, error) {

	filename := filepath.Base(zipFilePath)

	// If filename starts with "pftaps" then everything necessary is already known
	if strings.HasPrefix(filename, "pftaps") {
		return &types.USPTGoMetadata{
			DocumentType: "grant",
			OriginZip: types.OriginZip{
				ZipPath:       zipFilePath,
				ZipName:       filename,
				ZipEntryExt:   ".txt",
				Schema:        "aps",
				SchemaVersion: 0,
			},
		}, nil
	}

	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		// ! Implement error constructor and pass to channel
		return nil, err
	}
	defer zipReader.Close()

	type SchemaReference struct {
		Pattern       []byte
		DocumentType  string
		FileExtension string
		Schema        string
		SchemaVersion int8
	}

	// Define the patterns to match
	schemaReferences := []SchemaReference{
		{Pattern: []byte("us-patent-grant-v47-2022-02-17.dtd"), SchemaVersion: 47, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v47-2022-02-17.dtd"},
		{Pattern: []byte("us-patent-grant-v46-2021-08-30.dtd"), SchemaVersion: 46, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v46-2021-08-30.dtd"},
		{Pattern: []byte("us-patent-grant-v45-2014-04-03.dtd"), SchemaVersion: 45, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v45-2014-04-03.dtd"},
		{Pattern: []byte("us-patent-grant-v44-2013-05-16.dtd"), SchemaVersion: 44, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v44-2013-05-16.dtd"},
		{Pattern: []byte("us-patent-grant-v43-2012-12-04.dtd"), SchemaVersion: 43, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v43-2012-12-04.dtd"},
		{Pattern: []byte("us-patent-grant-v42-2006-08-23.dtd"), SchemaVersion: 42, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v42-2006-08-23.dtd"},
		{Pattern: []byte("us-patent-grant-v41-2005-08-25.dtd"), SchemaVersion: 41, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v41-2005-08-25.dtd"},
		{Pattern: []byte("us-patent-grant-v40-2004-12-02.dtd"), SchemaVersion: 40, DocumentType: "grant", FileExtension: ".xml", Schema: "us-patent-grant-v40-2004-12-02.dtd"},
		{Pattern: []byte("ST32-US-Grant-025xml.dtd"), SchemaVersion: 25, DocumentType: "grant", FileExtension: ".xml", Schema: "ST32-US-Grant-025xml.dtd"},
		{Pattern: []byte("us-patent-application-v46-2022-02-17.dtd"), SchemaVersion: 46, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v46-2022-02-17.dtd"},
		{Pattern: []byte("us-patent-application-v45-2021-08-30.dtd"), SchemaVersion: 45, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v45-2021-08-30.dtd"},
		{Pattern: []byte("us-patent-application-v44-2014-04-03.dtd"), SchemaVersion: 44, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v44-2014-04-03.dtd"},
		{Pattern: []byte("us-patent-application-v43-2012-12-04.dtd"), SchemaVersion: 43, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v43-2013-05-16.dtd"},
		{Pattern: []byte("us-patent-application-v42-2006-08-23.dtd"), SchemaVersion: 42, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v42-2012-12-04.dtd"},
		{Pattern: []byte("us-patent-application-v41-2005-08-25.dtd"), SchemaVersion: 41, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v41-2006-08-23.dtd"},
		{Pattern: []byte("us-patent-application-v40-2004-12-02.dtd"), SchemaVersion: 40, DocumentType: "application", FileExtension: ".xml", Schema: "us-patent-application-v40-2005-08-25.dtd"},
		{Pattern: []byte("pap-v16-2002-01-01.dtd"), SchemaVersion: 16, DocumentType: "application", FileExtension: ".xml", Schema: "pap-v16-2002-01-01.dtd"},
		{Pattern: []byte("pap-v15-2001-01-31.dtd"), SchemaVersion: 15, DocumentType: "application", FileExtension: ".xml", Schema: "pap-v15-2001-01-31.dtd"},
	}

	for _, zipEntry := range zipReader.File {
		// Skip directories
		if zipEntry.FileInfo().IsDir() {
			continue
		}

		// Only read .xml files within zip, ignoring .sgm files when they occasionally exist
		if strings.HasSuffix(strings.ToLower(zipEntry.Name), ".xml") {
			fileReader, err := zipEntry.Open()
			if err != nil {
				// ! Implement error constructor and pass to channel
				return nil, err // If it can't be opened, it must be skipped
			}
			defer fileReader.Close()

			// Read the first 2048 bytes of the file - more than enough to capture prolog
			buffer := make([]byte, 2048)
			_, err = fileReader.Read(buffer)
			if err != nil && err != io.EOF {
				// ! Implement error constructor and pass to channel
				return nil, err // If it can't be read, it must be skipped
			}

			// Iterate over the schema patterns to find a match
			for _, pattern := range schemaReferences {
				if bytes.Contains(buffer, pattern.Pattern) {
					return &types.USPTGoMetadata{
						DocumentType: pattern.DocumentType,
						OriginZip: types.OriginZip{
							ZipPath:       zipFilePath,
							ZipName:       filename,
							ZipEntryExt:   pattern.FileExtension,
							Schema:        pattern.Schema,
							SchemaVersion: pattern.SchemaVersion,
						},
					}, nil
				}
			}

		}

	}
	return nil, errors.New("failed to match a schema to zip")
}
