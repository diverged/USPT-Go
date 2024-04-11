package types

type USPTGoConfig struct {
	InputPath         string // Path to the input zip file
	ReturnRawSplitDoc bool   // Optional - return the raw split XML document in addition to the parsed document.  True by default.  False will save memory.
	Logger            Logger // Optional - provide a logger interface
}

// Logger defines a simple interface for logging within the parser.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

/*
USPTGoError enables centralized reporting of problems encounterd with specific files, especially skipped files.
Example Report: "%Type file skipped [%Name] due to error encountered while %Whence.\n Error: %v\n\n"
*/
type USPTGoError struct {
	Err     error  // The error encountered
	Skipped bool   // Whether the file was skipped
	Name    string // Zip name, Index within Zip, Document ID, etc.
	Whence  string // verb phrase, e.g. "opening the file", "reading the file", etc.
	Type    string // Zip, Part of Zip, Patent Doc, etc.
	ZipInfo OriginZip
}

func (e *USPTGoError) Error() string {
	return e.Err.Error()
}

// USPTGoDoc is the object returned via docChan
type USPTGoDoc struct {
	USPTGoMetadata USPTGoMetadata
	RawSplitDoc    []byte // Entire XML document as represented in the originating bulk file
	Patent         Patent
	Trademark      Trademark
}

// Trademark
type Trademark struct {
	RawSplitDoc []byte // Entire XML document as represented in the originating bulk file
}

// USPT-Go generated metadata
type USPTGoMetadata struct {
	DocumentType string
	OriginZip    OriginZip
}

// Type OriginZip contains metadata about the origin of a USPTGoDoc
type OriginZip struct {
	ZipPath       string
	ZipName       string
	ZipEntryExt   string
	Schema        string
	SchemaVersion int8
	IndexInZip    int
	IndexName     string
}
