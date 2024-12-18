package parser

import (
	"bytes"
	"github.com/russross/blackfriday/v2"
	"strings"
)

// ConvertMarkdownToHTML processes Markdown inside HTML tags and converts it to HTML.
func ConvertMarkdownToHTML(input string) (string, error) {
	// Step 1: Remove HTML comments from the input.
	inputWithoutComments := removeHTMLComments(input)

	// Step 2: Split the content into HTML and Markdown parts.
	parts := splitHTMLMarkdown(inputWithoutComments)

	// Step 3: Process each part and convert Markdown to HTML while preserving HTML parts.
	var output bytes.Buffer
	for _, part := range parts {
		if isHTML(part) {
			// Append HTML parts directly.
			output.WriteString(part)
		} else {
			// Convert Markdown parts to HTML using blackfriday.
			processedMarkdown := blackfriday.Run([]byte(part))
			output.Write(processedMarkdown)
		}
	}

	return output.String(), nil
}

// splitHTMLMarkdown splits the input string into alternating HTML and Markdown parts.
func splitHTMLMarkdown(input string) []string {
	var parts []string
	var current strings.Builder
	inTag := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		// Detect the start of an HTML tag.
		if char == '<' && !inTag {
			// Add the current Markdown content if any.
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			inTag = true
		}

		// Detect the end of an HTML tag.
		if char == '>' && inTag {
			inTag = false
		}

		// Append the current character to the current part.
		current.WriteByte(char)

		// If we've exited an HTML tag, save it as a separate part.
		if !inTag && char == '>' {
			parts = append(parts, current.String())
			current.Reset()
		}
	}

	// Add any remaining content as a Markdown part.
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// isHTML checks if a string is an HTML tag or element.
func isHTML(content string) bool {
	content = strings.TrimSpace(content)
	return len(content) > 1 && content[0] == '<' && content[len(content)-1] == '>'
}

// removeHTMLComments removes all HTML comments from the input string.
func removeHTMLComments(input string) string {
	var result strings.Builder
	inComment := false

	for i := 0; i < len(input); i++ {
		// Detect the start of an HTML comment.
		if i+3 < len(input) && input[i:i+4] == "<!--" {
			inComment = true
			i += 3 // Skip the opening comment sequence.
			continue
		}

		// Detect the end of an HTML comment.
		if inComment && i+2 < len(input) && input[i:i+3] == "-->" {
			inComment = false
			i += 2 // Skip the closing comment sequence.
			continue
		}

		// If not in a comment, write the character to the result.
		if !inComment {
			result.WriteByte(input[i])
		}
	}

	return result.String()
}
