package transformtext

import (
	"testing"
)

func TestTableXmlToHtml(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name: "Single table",
			input: []byte(`
                <description>
                    <table>
                        <title>Sample Table</title>
                        <tgroup cols="2">
                            <thead>
                                <row>
                                    <entry>Header 1</entry>
                                    <entry>Header 2</entry>
                                </row>
                            </thead>
                            <tbody>
                                <row>
                                    <entry>Row 1, Cell 1</entry>
                                    <entry>Row 1, Cell 2</entry>
                                </row>
                            </tbody>
                        </tgroup>
                    </table>
                </description>
            `),
			expected: []byte(`<description><table><caption>Sample Table</caption><colgroup><col span="1" align=""><col span="1" align=""></colgroup><thead><tr><th>Header 1</th><th>Header 2</th></tr></thead><tbody><tr><td>Row 1, Cell 1</td><td>Row 1, Cell 2</td></tr></tbody></table></description>`),
		},
		{
			name: "Multiple tables",
			input: []byte(`
                <description>
                    <table>
                        <title>Table 1</title>
                        <tgroup cols="2">
                            <tbody>
                                <row>
                                    <entry>Table 1, Row 1, Cell 1</entry>
                                    <entry>Table 1, Row 1, Cell 2</entry>
                                </row>
                            </tbody>
                        </tgroup>
                    </table>
                    <table>
                        <title>Table 2</title>
                        <tgroup cols="2">
                            <tbody>
                                <row>
                                    <entry>Table 2, Row 1, Cell 1</entry>
                                    <entry>Table 2, Row 1, Cell 2</entry>
                                </row>
                            </tbody>
                        </tgroup>
                    </table>
                </description>
            `),
			expected: []byte(`<description><table><caption>Table 1</caption><colgroup><col span="1" align=""><col span="1" align=""></colgroup><tbody><tr><td>Table 1, Row 1, Cell 1</td><td>Table 1, Row 1, Cell 2</td></tr></tbody></table><table><caption>Table 2</caption><colgroup><col span="1" align=""><col span="1" align=""></colgroup><tbody><tr><td>Table 2, Row 1, Cell 1</td><td>Table 2, Row 1, Cell 2</td></tr></tbody></table></description>`),
		},
		{
			name: "Self-closing tags",
			input: []byte(`
                <description>
                    <p>Paragraph with<br/>line break</p>
                    <table>
                        <title>Table</title>
                        <tgroup cols="2">
                            <tbody>
                                <row>
                                    <entry>Row 1, Cell 1</entry>
                                    <entry>Row 1, Cell 2</entry>
                                </row>
                            </tbody>
                        </tgroup>
                    </table>
                    <img src="example.jpg" alt="Example Image"/>
                </description>
            `),
			expected: []byte(`<description><p>Paragraph with<br/>line break</p><table><caption>Table</caption><colgroup><col span="1" align=""><col span="1" align=""></colgroup><tbody><tr><td>Row 1, Cell 1</td><td>Row 1, Cell 2</td></tr></tbody></table><img src="example.jpg" alt="Example Image"/></description>`),
		},
		{
			name: "No tables",
			input: []byte(`
                <description>
                    <p>Paragraph 1</p>
                    <p>Paragraph 2</p>
                </description>
            `),
			expected: []byte(`<description><p>Paragraph 1</p><p>Paragraph 2</p></description>`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TableXmlToHtml(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if string(result) != string(tc.expected) {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}
