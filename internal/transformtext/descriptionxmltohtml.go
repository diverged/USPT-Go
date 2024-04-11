package transformtext

import (
	"bytes"

	nethtml "golang.org/x/net/html"
)

func InnerXmlToHtml(descriptionXML []byte) (string, error) {

	// Remove all table elements
	descriptionXML, err := TableXmlToHtml(descriptionXML)
	if err != nil {
		return "", err
	}

	// Create a bytes.Reader from the []byte slice
	reader := bytes.NewReader(descriptionXML)

	// Parse the HTML content
	doc, err := nethtml.Parse(reader)
	if err != nil {
		return "", err
	}

	// Find the body element
	var body *nethtml.Node
	var f func(*nethtml.Node)
	f = func(n *nethtml.Node) {
		if n.Type == nethtml.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Traverse the HTML nodes and modify them as needed
	var processElement func(*nethtml.Node)
	processElement = func(n *nethtml.Node) {

		if n.Data == "table" || n.Data == "tables" {
			return // This skips the current node and all its children
		}

		switch n.Data {
		case "heading":
			// Change the node type to a header tag based on the level attribute
			level := "0" // Default level
			for i, a := range n.Attr {
				if a.Key == "level" || a.Key == "lvl" {
					level = a.Val
					n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
					break
				}
			}
			n.Data = "h2"
			class := "level-" + level
			n.Attr = append(n.Attr, nethtml.Attribute{Key: "class", Val: class})
		case "p":
			// Handle paragraphs with "h-" prefix in the id attribute
			id := ""
			for _, a := range n.Attr {
				if a.Key == "id" {
					id = a.Val
					break
				}
			}
			if id != "" && id[:2] == "h-" {
				level := id[2:]
				n.Data = "h4"
				class := "level-" + level
				n.Attr = append(n.Attr, nethtml.Attribute{Key: "class", Val: class})
			}
		case "ul", "ol":
			// Preserve the id attribute
			var idAttr string
			for i, a := range n.Attr {
				if a.Key == "id" {
					idAttr = a.Val
					n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
					break
				}
			}

			// Handle unordered lists with specific list styles
			if n.Data == "ul" {
				listStyle := ""
				for _, a := range n.Attr {
					if a.Key == "list-style" {
						listStyle = a.Val
						break
					}
				}
				switch listStyle {
				case "none":
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}, {Key: "style", Val: "list-style-type:none"}}
				case "bullet":
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}, {Key: "style", Val: "list-style-type:disc"}}
				case "dash":
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}, {Key: "class", Val: "ul-dash"}}
				default:
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}}
				}
			}

			// Handle ordered lists with specific list styles
			if n.Data == "ol" {
				listStyle := ""
				for _, a := range n.Attr {
					if a.Key == "style" || a.Key == "ol-style" {
						listStyle = a.Val
						break
					}
				}
				if listStyle != "" {
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}, {Key: "type", Val: listStyle}}
				} else {
					n.Attr = []nethtml.Attribute{{Key: "id", Val: idAttr}}
				}
			}
		case "br":
			// Remove any extra text content for self-closing elements
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == nethtml.TextNode {
					c.Data = ""
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processElement(c)
		}
	}

	// Process the elements within the body
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		processElement(c)
	}

	// Render the modified document back to a string
	var buf bytes.Buffer
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == nethtml.TextNode {
			buf.WriteString(c.Data)
		} else if c.Type == nethtml.ElementNode {
			if err := nethtml.Render(&buf, c); err != nil {
				return "", err
			}
		}
	}

	return buf.String(), nil
}
