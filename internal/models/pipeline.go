package models

type SplitXMLDoc struct {
	Content    []byte // Content of single document split from bulk file
	IndexInZip map[string]string
}
