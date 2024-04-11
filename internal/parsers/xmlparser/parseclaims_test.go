package xmlparser

import (
	"testing"

	"github.com/diverged/uspt-go/internal/models"
)

// mockLogger implements the Logger interface for testing purposes.
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) {}

// TestParseStructuredClaims tests the ParseStructuredClaims function for correct parsing of XML claims.
func TestParseStructuredClaims(t *testing.T) {
	// Sample XML input
	xmlInput := []byte(`
	<claims>
		<claim id="1">
			<claim-text>Claim 1 text.</claim-text>
		</claim>
		<claim id="2">
			<claim-text>Claim 2 text with <claim-ref idref="1">claim 1</claim-ref>.</claim-text>
		</claim>
	</claims>
	`)

	// Expected output
	expectedClaims := []*models.Claim{
		{
			ID:   "1",
			Type: "INDEPENDENT",
			Text: []string{"Claim 1 text."},
			ClaimTree: models.ClaimTree{
				ParentCount:    0,
				ChildCount:     1,
				ClaimTreeLevel: 0,
			},
		},
		{
			ID:   "2",
			Type: "DEPENDENT",
			Text: []string{"Claim 2 text with claim 1."},
			ClaimTree: models.ClaimTree{
				ParentIds:      []string{"1"},
				ParentCount:    1,
				ChildCount:     0,
				ClaimTreeLevel: 1,
			},
		},
	}

	log := &mockLogger{}
	claims, err := ParseStructuredClaims(xmlInput, log)
	if err != nil {
		t.Errorf("ParseStructuredClaims returned an error: %v", err)
	}

	// Verify the length of the claims slice
	if len(claims) != len(expectedClaims) {
		t.Fatalf("Expected %d claims, got %d", len(expectedClaims), len(claims))
	}

	// Verify the content of each claim
	for i, claim := range claims {
		if claim.ID != expectedClaims[i].ID || claim.Type != expectedClaims[i].Type || len(claim.Text) != len(expectedClaims[i].Text) {
			t.Errorf("Claim %d does not match expected output", i+1)
		}
	}
}
