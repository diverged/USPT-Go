## USPTGo - USPTO Bulk Data Processing in Go

A Go package which accepts U.S. Patent and Trademark Office (USPTO) [bulk data zip files](https://bulkdata.uspto.gov/), and returns standardized objects of structured, formatted patent contents.

For a standalone tool implementation of this package, see [USPTO-Bulk-Data-Tool](github.com/diverged/uspto-bulk-data-tool/).

At this time, the USPTGo package supports the following USPTO bulk data products:
- **Patent Grant Full Text Data (No Images) (2004 - Present)**
- **Patent Application Full Text Data (No Images) (2004 - Present)**

### Usage

```go
func USPTGo(cfg *types.USPTGoConfig) (<-chan *types.USPTGoDoc, <-chan error, error)
```

Process a bulk data zip by passing an instance of USPTGoConfig to the USPTGo function, which returns two buffered channels, and an error.

```go
type USPTGoConfig struct {
	InputPath         string // Path to the input zip file
	ReturnRawSplitDoc bool   // Optional - returns the raw split XML document in addition to the parsed document.  True by default.  False will save memory.
	Logger            Logger // Optional - provide a logging interface
}
```

The first channel returned contains individual documents from the inputted zip file:

```go
type USPTGoDoc struct {
	USPTGoMetadata USPTGoMetadata
	RawSplitDoc    []byte // Entire XML document as represented in the originating bulk file
	Patent         Patent
	Trademark      Trademark 
}

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
```

The second channel contains errors encountered, including information like whether or not a document was skipped.


```go
type USPTGoError struct {
	Err     error  // The error encountered
	Skipped bool   // Whether the file was skipped
	Name    string // Zip name, Index within Zip, Document ID, etc.
	Whence  string // verb phrase, e.g. "opening the file", "reading the file", etc.
	Type    string // Zip, Part of Zip, Patent Doc, etc.
	ZipInfo OriginZip
}
```

### Example

Minimal example:

```go
package main

import (
    "github.com/diverged/uspt-go"
    "github.com/diverged/uspt-go/types"
)

func main() {
    cfg := &types.USPTGoConfig{
        // Initialize your config
    }

    docChan, errChan, err := usptgo.USPTGo(cfg)
    if err != nil {
        // Handle initialization error
    }

    // Example of how to use the returned channels
    for doc := range docChan {
        // Process each document
    }

    for err := range errChan {
        // Handle each error
    }
}
```

For a more complete example of how to make use of this package, see [USPTO-Bulk-Data-Tool](github.com/diverged/uspto-bulk-data-tool/).

### License

[MIT](https://github.com/diverged/USPT-Go/blob/main/LICENSE)