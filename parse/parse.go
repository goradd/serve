package parse

import (
	"bytes"

	"golang.org/x/net/html"
)

type GoraddElement struct {
	Tag   string // e.g. "div"
	ID    string // id attribute value (if any)
	Start int    // byte offset of the '<' of the opening tag
	End   int    // byte offset just *after* the '>' of the closing (or self-closing) tag
}

// void elements in HTML (no closing tag)
// For these, the element ends at the end of its start tag.
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// stackElem tracks open tags while tokenizing
type stackElem struct {
	tag       string
	startPos  int
	hasGoradd bool
	id        string
}

// findGoraddElements scans the HTML source and returns all elements
// that have a "data-goradd" attribute, along with their positions.
func findGoraddElements(src []byte) ([]GoraddElement, error) {
	z := html.NewTokenizer(bytes.NewReader(src))

	var (
		offset int // running byte offset through src
		stack  []stackElem
		result []GoraddElement
	)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			// ErrorToken is returned at EOF; any other error is inside z.Err()
			if z.Err() != nil && z.Err().Error() != "EOF" {
				return nil, z.Err()
			}
			break
		}

		raw := z.Raw() // raw bytes of this token, exactly as in the input
		tokenStart := offset
		tokenEnd := offset + len(raw)
		offset = tokenEnd

		tok := z.Token()
		tagName := tok.Data

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			// inspect attributes
			hasGoradd := false
			idVal := ""
			for _, attr := range tok.Attr {
				if attr.Key == "data-goradd" {
					hasGoradd = true
				}
				if attr.Key == "id" {
					idVal = attr.Val
				}
			}

			// SelfClosingTagToken or known void elements end right here.
			isSelfClosing := tt == html.SelfClosingTagToken || voidElements[tagName]

			if isSelfClosing {
				if hasGoradd {
					result = append(result, GoraddElement{
						Tag:   tagName,
						ID:    idVal,
						Start: tokenStart,
						End:   tokenEnd,
					})
				}
				// no need to push onto stack
				continue
			}

			// Normal start tag: push onto stack
			stack = append(stack, stackElem{
				tag:       tagName,
				startPos:  tokenStart,
				hasGoradd: hasGoradd,
				id:        idVal,
			})

		case html.EndTagToken:
			// Find the matching start tag in the stack (walk backwards)
			i := len(stack) - 1
			for i >= 0 && stack[i].tag != tagName {
				i--
			}
			if i < 0 {
				// No matching start tag found (malformed HTML); ignore
				continue
			}

			// Pop the matched element
			elem := stack[i]
			stack = stack[:i]

			if elem.hasGoradd {
				result = append(result, GoraddElement{
					Tag:   elem.tag,
					ID:    elem.id,
					Start: elem.startPos,
					End:   tokenEnd,
				})
			}
		}
	}

	return result, nil
}
