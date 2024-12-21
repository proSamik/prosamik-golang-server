package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// ConvertMarkdownToHTML converts Markdown to HTML using gomarkdown library
func ConvertMarkdownToHTML(input string) (string, error) {
	// Validate input
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	// Preprocess the input to handle nested lists and special formatting
	input = preprocessMarkdown(input)

	// Create a Markdown parser with comprehensive extensions
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.Strikethrough |
		parser.Footnotes |
		parser.HeadingIDs |
		parser.OrderedListStart |
		parser.NoIntraEmphasis // Prevent unwanted emphasis

	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer with comprehensive options
	htmlFlags := html.CommonFlags |
		html.HrefTargetBlank

	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	// Convert Markdown to HTML
	md := []byte(input)
	htmlContent := markdown.ToHTML(md, p, renderer)

	return string(htmlContent), nil
}

// preprocessMarkdown handles special Markdown formatting cases
func preprocessMarkdown(input string) string {
	// Replace Windows-style line breaks with Unix-style
	input = strings.ReplaceAll(input, "\r\n", "\n")

	// Regex for detecting list items with various formatting
	listItemRegex := regexp.MustCompile(`^([ \t]*)([-*]|\d+\.)\s*(.*)$`)

	// Split input into lines
	lines := strings.Split(input, "\n")
	var processedLines []string
	var inListBlock bool

	for _, line := range lines {
		matches := listItemRegex.FindStringSubmatch(line)

		if matches != nil {
			// This is a list item
			indent := matches[1]
			marker := matches[2]
			content := matches[3]

			// Ensure we start a list block if not already in one
			if !inListBlock {
				inListBlock = true
				// Add an empty line before the list to ensure proper list parsing
				processedLines = append(processedLines, "")
			}

			// Reconstruct the list item
			processedLines = append(processedLines, indent+marker+" "+content)
		} else {
			// Non-list line
			if inListBlock && strings.TrimSpace(line) != "" {
				// If we were in a list block and this is not an empty line,
				// it means the list has ended
				inListBlock = false
				// Add an empty line to separate lists
				processedLines = append(processedLines, "")
			}

			processedLines = append(processedLines, line)
		}
	}

	// Join the processed lines
	return strings.Join(processedLines, "\n")
}

// RemoveHTMLComments removes HTML comments from the input string
func RemoveHTMLComments(input string) string {
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	return commentRegex.ReplaceAllString(input, "")
}

// ProcessHTMLWithMarkdown processes Markdown within HTML, converting Markdown parts to HTML
func ProcessHTMLWithMarkdown(input string) (string, error) {
	// Remove comments first
	input = RemoveHTMLComments(input)

	// Regular expression to split HTML and Markdown
	splitRegex := regexp.MustCompile(`(<[^>]+>|</[^>]+>)`)
	parts := splitRegex.Split(input, -1)

	var processedHTML bytes.Buffer
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)

		// If it's an HTML tag, write it directly
		if splitRegex.MatchString(trimmedPart) {
			processedHTML.WriteString(trimmedPart)
			continue
		}

		// If it's Markdown content, convert to HTML
		if trimmedPart != "" {
			htmlContent, err := ConvertMarkdownToHTML(trimmedPart)
			if err != nil {
				return "", fmt.Errorf("error converting Markdown to HTML: %v", err)
			}
			processedHTML.WriteString(htmlContent)
		}
	}

	return processedHTML.String(), nil
}
