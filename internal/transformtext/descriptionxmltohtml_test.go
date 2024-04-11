package transformtext

import (
	"testing"
)

func TestInnerXmlToHtml(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{

		{
			name:     "Simple heading",
			input:    `<heading id="h1">Test Heading</heading>`,
			expected: `<h2 id="h1" class="level-0">Test Heading</h2>`,
		},
		{
			name:     "Heading with level attribute",
			input:    `<heading id="h2" level="2">Subheading</heading>`,
			expected: `<h2 id="h2" class="level-2">Subheading</h2>`,
		},
		{
			name:     "Heading with lvl attribute",
			input:    `<heading id="h3" lvl="3">Sub-subheading</heading>`,
			expected: `<h2 id="h3" class="level-3">Sub-subheading</h2>`,
		},
		{
			name:     "Paragraph",
			input:    `<p>This is a paragraph.</p>`,
			expected: `<p>This is a paragraph.</p>`,
		},
		{
			name:     "Paragraph with h- id",
			input:    `<p id="h-1">This is a paragraph with h- id.</p>`,
			expected: `<h4 id="h-1" class="level-1">This is a paragraph with h- id.</h4>`,
			//expected: `<h4 id="h-1" class="level-0">This is a paragraph with h- id.</h4>`,
		},
		{
			name:     "Unordered list with default style",
			input:    `<ul id="list1"><li>Item 1</li><li>Item 2</li></ul>`,
			expected: `<ul id="list1"><li>Item 1</li><li>Item 2</li></ul>`,
		},
		{
			name:     "Unordered list with none style",
			input:    `<ul id="list2" list-style="none"><li>Item 1</li><li>Item 2</li></ul>`,
			expected: `<ul id="list2" style="list-style-type:none"><li>Item 1</li><li>Item 2</li></ul>`,
		},
		{
			name:     "Unordered list with bullet style",
			input:    `<ul id="list3" list-style="bullet"><li>Item 1</li><li>Item 2</li></ul>`,
			expected: `<ul id="list3" style="list-style-type:disc"><li>Item 1</li><li>Item 2</li></ul>`,
		},
		{
			name:     "Unordered list with dash style",
			input:    `<ul id="list4" list-style="dash"><li>Item 1</li><li>Item 2</li></ul>`,
			expected: `<ul id="list4" class="ul-dash"><li>Item 1</li><li>Item 2</li></ul>`,
		},
		{
			name:     "Ordered list with default style",
			input:    `<ol id="list5"><li>Item 1</li><li>Item 2</li></ol>`,
			expected: `<ol id="list5"><li>Item 1</li><li>Item 2</li></ol>`,
		},
		{
			name:     "Ordered list with style attribute",
			input:    `<ol id="list6" style="a"><li>Item 1</li><li>Item 2</li></ol>`,
			expected: `<ol id="list6" type="a"><li>Item 1</li><li>Item 2</li></ol>`,
		},
		{
			name:     "Nested elements",
			input:    `<p>This is a <b>paragraph</b> with <i>nested</i> elements.</p>`,
			expected: `<p>This is a <b>paragraph</b> with <i>nested</i> elements.</p>`,
		},
		{
			name:     "Self-closing elements",
			input:    `<p>This is a paragraph with a <br/> self-closing element.</p>`,
			expected: `<p>This is a paragraph with a <br/> self-closing element.</p>`,
		},

		{
			name:     "Nested lists",
			input:    `<ul><li>Item 1<ul><li>Nested Item 1</li><li>Nested Item 2</li></ul></li><li>Item 2</li></ul>`,
			expected: `<ul id=""><li>Item 1<ul id=""><li>Nested Item 1</li><li>Nested Item 2</li></ul></li><li>Item 2</li></ul>`,
		},
		{
			name:     "Empty elements",
			input:    `<p></p><ul></ul><ol></ol>`,
			expected: `<p></p><ul id=""></ul><ol id=""></ol>`,
		},
		{
			name:     "Special characters",
			input:    `<p>This is a paragraph with &lt;, &gt;, &amp;, &#34;, and &#39; characters.</p>`,
			expected: `<p>This is a paragraph with &lt;, &gt;, &amp;, &#34;, and &#39; characters.</p>`,
		},
		{
			name:     "Whitespace handling",
			input:    `<p>This is a paragraph with    extra   whitespace.</p>`,
			expected: `<p>This is a paragraph with    extra   whitespace.</p>`,
		},
		{
			name:     "Unknown elements",
			input:    `<unknown>This is an unknown element.</unknown>`,
			expected: `<unknown>This is an unknown element.</unknown>`,
		},
		{
			name:     "Simple table",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry>Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with thead",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><thead><row><entry>Header 1</entry><entry>Header 2</entry></row></thead><tbody><row><entry>Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with pgwide attribute",
			input:    `<table pgwide="1"><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry>Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with frame attribute",
			input:    `<table frame="all"><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry>Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with colspec attributes",
			input:    `<table><tgroup><colspec colname="col1" colwidth="50" align="left"/><colspec colname="col2" colwidth="100" align="center"/><tbody><row><entry>Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with entry attributes",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry valign="top" align="left">Cell 1</entry><entry valign="bottom" align="right">Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with morerows attribute",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry morerows="1">Cell 1</entry><entry>Cell 2</entry></row><row><entry>Cell 3</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with namest and nameend attributes",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><colspec colname="col3"/><tbody><row><entry namest="col1" nameend="col2">Cell 1</entry><entry>Cell 2</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
		{
			name:     "Table with rowsep and colsep attributes",
			input:    `<table><tgroup><colspec colname="col1"/><colspec colname="col2"/><tbody><row><entry rowsep="1" colsep="1">Cell 1</entry><entry rowsep="1" colsep="0">Cell 2</entry></row><row><entry rowsep="0" colsep="1">Cell 3</entry><entry rowsep="0" colsep="0">Cell 4</entry></row></tbody></tgroup></table>`,
			expected: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := InnerXmlToHtml([]byte(tc.input))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected: %s\nGot: %s", tc.expected, result)
			}
		})
	}
}
