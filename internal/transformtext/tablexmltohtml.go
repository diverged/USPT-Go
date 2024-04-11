package transformtext

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Table struct {
	XMLName xml.Name `xml:"table"`
	Title   string   `xml:"title"`
	TGroups []TGroup `xml:"tgroup"`
}

type TGroup struct {
	XMLName xml.Name  `xml:"tgroup"`
	Cols    int       `xml:"cols,attr"`
	ColSpec []ColSpec `xml:"colspec"`
	THead   *THead    `xml:"thead"`
	TBody   []TBody   `xml:"tbody"`
}

type ColSpec struct {
	ColNum  int    `xml:"colnum,attr"`
	ColName string `xml:"colname,attr"`
	Align   string `xml:"align,attr"`
}

type THead struct {
	XMLName xml.Name `xml:"thead"`
	Rows    []Row    `xml:"row"`
}

type TBody struct {
	XMLName xml.Name `xml:"tbody"`
	Rows    []Row    `xml:"row"`
}

type Row struct {
	XMLName xml.Name `xml:"row"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName  xml.Name `xml:"entry"`
	NameSt   string   `xml:"namest,attr"`
	NameEnd  string   `xml:"nameend,attr"`
	MoreRows int      `xml:"morerows,attr"`
	Content  string   `xml:",innerxml"`
}

func (t *Table) translateTableXmlToHtml() string {
	var sb strings.Builder

	sb.WriteString("<table>")
	if t.Title != "" {
		sb.WriteString("<caption>" + t.Title + "</caption>")
	}

	for _, tgroup := range t.TGroups {
		sb.WriteString("<colgroup>")
		if len(tgroup.ColSpec) == 0 {
			// If no colspec is provided, generate default col tags
			for i := 0; i < tgroup.Cols; i++ {
				sb.WriteString("<col span=\"1\" align=\"\">")
			}
		} else {
			for _, colspec := range tgroup.ColSpec {
				sb.WriteString(fmt.Sprintf("<col span=\"1\" align=\"%s\">", colspec.Align))
			}
		}
		sb.WriteString("</colgroup>")

		if tgroup.THead != nil {
			sb.WriteString("<thead>")
			for _, row := range tgroup.THead.Rows {
				sb.WriteString("<tr>")
				for _, entry := range row.Entries {
					sb.WriteString("<th>" + entry.Content + "</th>")
				}
				sb.WriteString("</tr>")
			}
			sb.WriteString("</thead>")
		}

		for _, tbody := range tgroup.TBody {
			sb.WriteString("<tbody>")
			for _, row := range tbody.Rows {
				sb.WriteString("<tr>")
				for _, entry := range row.Entries {
					if entry.NameSt != "" && entry.NameEnd != "" {
						sb.WriteString(fmt.Sprintf("<td colspan=\"%d\">%s</td>", tgroup.Cols, entry.Content))
					} else if entry.MoreRows > 0 {
						sb.WriteString(fmt.Sprintf("<td rowspan=\"%d\">%s</td>", entry.MoreRows+1, entry.Content))
					} else {
						sb.WriteString("<td>" + entry.Content + "</td>")
					}
				}
				sb.WriteString("</tr>")
			}
			sb.WriteString("</tbody>")
		}
	}

	sb.WriteString("</table>")
	return sb.String()
}

/* func (t *Table) translateTableXmlToHtml() string {
	var sb strings.Builder

	sb.WriteString("<table>")
	if t.Title != "" {
		sb.WriteString("<caption>" + t.Title + "</caption>")
	}

	for _, tgroup := range t.TGroups {
		sb.WriteString("<colgroup>")
		for _, colspec := range tgroup.ColSpec {
			sb.WriteString(fmt.Sprintf("<col span=\"1\" align=\"%s\">", colspec.Align))
		}
		sb.WriteString("</colgroup>")

		if tgroup.THead != nil {
			sb.WriteString("<thead>")
			for _, row := range tgroup.THead.Rows {
				sb.WriteString("<tr>")
				for _, entry := range row.Entries {
					sb.WriteString("<th>" + entry.Content + "</th>")
				}
				sb.WriteString("</tr>")
			}
			sb.WriteString("</thead>")
		}

		for _, tbody := range tgroup.TBody {
			sb.WriteString("<tbody>")
			for _, row := range tbody.Rows {
				sb.WriteString("<tr>")
				for _, entry := range row.Entries {
					if entry.NameSt != "" && entry.NameEnd != "" {
						sb.WriteString(fmt.Sprintf("<td colspan=\"%d\">%s</td>", tgroup.Cols, entry.Content))
					} else if entry.MoreRows > 0 {
						sb.WriteString(fmt.Sprintf("<td rowspan=\"%d\">%s</td>", entry.MoreRows+1, entry.Content))
					} else {
						sb.WriteString("<td>" + entry.Content + "</td>")
					}
				}
				sb.WriteString("</tr>")
			}
			sb.WriteString("</tbody>")
		}
	}

	sb.WriteString("</table>")
	return sb.String()
} */

func TableXmlToHtml(descriptionXML []byte) ([]byte, error) {
	xmlData := string(descriptionXML)

	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	var builder strings.Builder

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "table" {
				var table Table
				err := decoder.DecodeElement(&table, &token)
				if err != nil {
					fmt.Println("Error parsing table element:", err)
					continue
				}

				htmlTable := table.translateTableXmlToHtml()
				builder.WriteString(htmlTable)
			} else {
				builder.WriteString("<" + token.Name.Local)
				for _, attr := range token.Attr {
					builder.WriteString(" " + attr.Name.Local + "=\"" + attr.Value + "\"")
				}
				if token.Name.Local == "br" || token.Name.Local == "img" {
					builder.WriteString("/>")
				} else {
					builder.WriteString(">")
				}
			}
		case xml.EndElement:
			if token.Name.Local != "table" && token.Name.Local != "br" && token.Name.Local != "img" {
				builder.WriteString("</" + token.Name.Local + ">")
			}
		case xml.CharData:

			/* 			data := strings.TrimSpace(string(token))
			   			if data != "" {
			   				builder.WriteString(data)
			   			} */

			data := string(token)
			data = strings.ReplaceAll(data, "\n", "")
			data = strings.ReplaceAll(data, "\t", "")
			builder.WriteString(data)
		default:
			// Ignore other token types
		}
	}

	// transformedXML := strings.TrimSpace(builder.String())
	transformedXML := builder.String()
	return []byte(transformedXML), nil
}
