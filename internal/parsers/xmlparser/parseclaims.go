package xmlparser

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/diverged/uspt-go/internal/models"
	"github.com/diverged/uspt-go/types"
)

func ParseStructuredClaims(rawXmlClaims []byte, log types.Logger) ([]*models.Claim, error) {

	var xmlClaims []models.XMLClaim
	decoder := xml.NewDecoder(bytes.NewReader(rawXmlClaims))
	//decoder := xml.NewDecoder(strings.NewReader(xmlData))
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				// EOF is expected, so we can break without reporting an error.
				break
			}
			log.Error("Error within ParseStructuredClaims when decoding token:", err)
			break
		}

		if startElement, ok := token.(xml.StartElement); ok && startElement.Name.Local == "claim" {
			var xmlClaim models.XMLClaim
			err := decoder.DecodeElement(&xmlClaim, &startElement)
			if err != nil {
				// fmt.Println("Error unmarshalling claim:", err)
				log.Error("Error unmarshalling claim:", err)
				continue
			}
			xmlClaims = append(xmlClaims, xmlClaim)
		}
	}

	//var claims []*Claim
	var claims []*models.Claim

	for _, xmlClaim := range xmlClaims {
		claim := &models.Claim{
			ID:   xmlClaim.ID,
			Type: "INDEPENDENT",
			Text: []string{},
			ClaimTree: models.ClaimTree{
				ParentCount:    0,
				ChildCount:     0,
				ClaimTreeLevel: 0,
			},
		}

		// Process main claim text
		claimText := xmlClaim.ClaimText.Text
		claimText = strings.ReplaceAll(claimText, "<claim-text>", "")
		claimText = strings.ReplaceAll(claimText, "</claim-text>", "")
		claimText = removeClaimRefTags(claimText)
		claim.Text = append(claim.Text, claimText)

		// Process nested claim text elements
		for _, element := range xmlClaim.ClaimText.Elements {
			elementText := element.Text
			elementText = strings.ReplaceAll(elementText, "<claim-text>", "")
			elementText = strings.ReplaceAll(elementText, "</claim-text>", "")
			elementText = removeClaimRefTags(elementText)
			claim.Text = append(claim.Text, elementText)

			// Process claim references
			for _, claimRef := range element.ClaimRefs {
				claim.Type = "DEPENDENT"
				claim.ClaimTree.ParentIds = append(claim.ClaimTree.ParentIds, claimRef.IDRef)
				claim.ClaimTree.ParentCount++
			}
		}

		// Check if the claim has a parent reference
		if xmlClaim.ClaimText.ClaimRef.IDRef != "" {
			claim.Type = "DEPENDENT"
			claim.ClaimTree.ParentIds = append(claim.ClaimTree.ParentIds, xmlClaim.ClaimText.ClaimRef.IDRef)
			claim.ClaimTree.ParentCount++
		}

		claims = append(claims, claim)
	}

	// Build the claim tree
	for _, claim := range claims {
		for _, parentID := range claim.ClaimTree.ParentIds {
			for _, parent := range claims {
				if parent.ID == parentID {
					parent.ChildIds = append(parent.ChildIds, claim.ID)
					parent.ClaimTree.ChildCount++
					break
				}
			}
		}
	}

	// Calculate claim tree levels
	var calculateClaimTreeLevel func(*models.Claim, int)
	calculateClaimTreeLevel = func(claim *models.Claim, level int) {
		claim.ClaimTree.ClaimTreeLevel = level
		for _, childID := range claim.ChildIds {
			for _, child := range claims {
				if child.ID == childID {
					calculateClaimTreeLevel(child, level+1)
					break
				}
			}
		}
	}

	for _, claim := range claims {
		if claim.Type == "INDEPENDENT" {
			calculateClaimTreeLevel(claim, 0)
		}
	}

	/* 	// Print the claim data
	   	for _, claim := range claims {
	   		fmt.Printf("Claim ID: %s\n", claim.ID)
	   		fmt.Printf("Type: %s\n", claim.Type)
	   		fmt.Printf("Text:\n")
	   		for _, text := range claim.Text {
	   			fmt.Printf("  %s\n", text)
	   		}
	   		fmt.Printf("Parent IDs: %v\n", claim.ClaimTree.ParentIds)
	   		fmt.Printf("Parent Count: %d\n", claim.ClaimTree.ParentCount)
	   		fmt.Printf("Child IDs: %v\n", claim.ChildIds)
	   		fmt.Printf("Child Count: %d\n", claim.ClaimTree.ChildCount)
	   		fmt.Printf("Claim Tree Level: %d\n", claim.ClaimTree.ClaimTreeLevel)
	   		fmt.Println("---")
	   	} */
	return claims, nil
}

// Helper function to remove <claim-ref> tags
func removeClaimRefTags(text string) string {
	regex := regexp.MustCompile(`<claim-ref[^>]*>|</claim-ref>`)
	return regex.ReplaceAllString(text, "")
}
