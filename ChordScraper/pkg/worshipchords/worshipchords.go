package worshipchords

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Song represents a worship song with chords
type Song struct {
	Title   string
	Artist  string
	Key     string
	Content string
	URL     string
}

// Client handles requests to worshipchords.com
type Client struct {
	httpClient *http.Client
}

// New creates a new worshipchords client
func New() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// GetSongFromURL fetches and parses a song from worshipchords.com
func (c *Client) GetSongFromURL(url string) (*Song, error) {
	// Validate URL
	if !strings.Contains(url, "worshipchords.com") {
		return nil, fmt.Errorf("invalid URL: must be from worshipchords.com")
	}

	// Fetch the page
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the HTML
	song, err := c.parseHTML(string(body), url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return song, nil
}

// parseHTML extracts song information from the HTML content
func (c *Client) parseHTML(htmlContent string, url string) (*Song, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	song := &Song{URL: url}

	// Find the title
	song.Title = findMetaContent(doc, "og:title")
	if song.Title == "" {
		song.Title = findTitle(doc)
	}

	// Extract song name and artist from title
	// Format is usually "Song Name Chords - Artist Name"
	if strings.Contains(song.Title, " - ") {
		parts := strings.Split(song.Title, " - ")
		songName := strings.TrimSpace(strings.Replace(parts[0], " Chords", "", 1))
		artistName := strings.TrimSpace(parts[1])
		song.Title = songName
		song.Artist = artistName
	}

	// Find the chord content
	chordContent := findSongChords(doc)
	if chordContent == "" {
		return nil, fmt.Errorf("no chord content found")
	}

	// Extract key if present
	song.Key = extractKey(chordContent)

	// Format the content
	song.Content = formatChordContent(chordContent)

	return song, nil
}

// findSongChords finds the div with class="song-chords-content" and extracts the pre element
func findSongChords(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "song-chords-content") {
				// Found the div, now find the pre element
				return findPreContent(n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := findSongChords(c)
		if result != "" {
			return result
		}
	}

	return ""
}

// findPreContent finds pre element and extracts its text content
func findPreContent(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "pre" {
		return extractText(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := findPreContent(c)
		if result != "" {
			return result
		}
	}

	return ""
}

// extractText extracts all text content from a node, preserving structure
func extractText(n *html.Node) string {
	var sb strings.Builder

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		} else if node.Type == html.ElementNode {
			if node.Data == "br" {
				sb.WriteString("\n")
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return sb.String()
}

// findMetaContent finds meta tag content by property name
func findMetaContent(n *html.Node, property string) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var prop, content string
		for _, attr := range n.Attr {
			if attr.Key == "property" && attr.Val == property {
				prop = attr.Val
			}
			if attr.Key == "content" {
				content = attr.Val
			}
		}
		if prop == property && content != "" {
			return content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := findMetaContent(c, property)
		if result != "" {
			return result
		}
	}

	return ""
}

// findTitle finds the title element
func findTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := findTitle(c)
		if result != "" {
			return result
		}
	}

	return ""
}

// extractKey attempts to find the key in the content
func extractKey(content string) string {
	// Look for "data-key" attribute in the content
	keyRegex := regexp.MustCompile(`data-key="([^"]+)"`)
	matches := keyRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// formatChordContent cleans up and formats the chord content for OnSong
func formatChordContent(content string) string {
	// Split into lines
	lines := strings.Split(content, "\n")
	var formatted []string
	var lastWasEmpty bool
	var lastWasSectionHeader bool
	var lastWasChordLine bool

	// Regex patterns for chord detection and section headers
	// Match section headers like "Verse 1", "Chorus", "Bridge", etc.
	sectionPattern := regexp.MustCompile(`^(Intro|Verse\s*\d*|Chorus\s*\d*|Pre-Chorus|Bridge|Instrumental|Interlude|Turnaround|Outro|Tag|Ending|Solo|Break|Refrain|Coda|Hook|Vamp|Outro Chorus)\s*$`)

	for _, line := range lines {
		// Trim whitespace from the right
		cleaned := strings.TrimRight(line, " \t")

		// Skip empty lines at the beginning
		if len(formatted) == 0 && cleaned == "" {
			continue
		}

		// Handle empty lines
		if cleaned == "" {
			// Skip blank line right after section header
			if lastWasSectionHeader {
				continue
			}
			// Skip blank line right after chord line (chords should be directly above lyrics)
			if lastWasChordLine {
				continue
			}
			// Only allow single blank lines between content
			if !lastWasEmpty {
				formatted = append(formatted, "")
				lastWasEmpty = true
			}
			continue
		}
		lastWasEmpty = false

		// Check if this is a section header
		if sectionPattern.MatchString(cleaned) {
			// Add colon to section headers
			formatted = append(formatted, cleaned+":")
			lastWasSectionHeader = true
			lastWasChordLine = false
			continue
		}
		lastWasSectionHeader = false

		// Check if this line looks like a chord line
		// Apply chord wrapping to lines that are primarily chords
		if isChordLine(cleaned) {
			// Wrap each chord in brackets
			wrappedLine := wrapChordsInBrackets(cleaned)
			formatted = append(formatted, wrappedLine)
			lastWasChordLine = true
		} else {
			// Regular lyric line or other content
			formatted = append(formatted, cleaned)
			lastWasChordLine = false
		}
	}

	// Join lines back together
	result := strings.Join(formatted, "\n")

	// Final cleanup: ensure no more than one consecutive blank line
	result = regexp.MustCompile(`\n\n+`).ReplaceAllString(result, "\n\n")

	return strings.TrimSpace(result)
}

// wrapChordsInBrackets wraps chord symbols in brackets
func wrapChordsInBrackets(line string) string {
	// Pattern to match chord symbols more comprehensively
	// Matches: A, Am, A7, Amaj7, Asus4, A/C#, etc.
	chordPattern := regexp.MustCompile(`\b([A-G][#b]?(?:maj|min|m|sus|aug|dim|add)?\d*(?:/[A-G][#b]?)?)\b`)

	// Replace each chord with [chord]
	result := chordPattern.ReplaceAllString(line, "[$1]")

	return result
}

// containsChords checks if a line contains chord patterns
func containsChords(line string) bool {
	chordPattern := regexp.MustCompile(`\b([A-G][#b]?(?:maj|min|m|sus|aug|dim|add)?\d*(?:/[A-G][#b]?)?)\b`)
	return chordPattern.MatchString(line)
}

// isChordLine determines if a line is primarily chords vs lyrics
func isChordLine(line string) bool {
	// If line is empty, not a chord line
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}

	// Split line into tokens (words)
	tokens := strings.Fields(trimmed)
	if len(tokens) == 0 {
		return false
	}

	// Pattern for valid chord tokens
	chordPattern := regexp.MustCompile(`^[A-G][#b]?(?:maj|min|m|sus|aug|dim|add)?\d*(?:/[A-G][#b]?)?$`)
	chordCount := 0

	for _, token := range tokens {
		if chordPattern.MatchString(token) {
			chordCount++
		}
	}

	// If ALL tokens are chords, definitely a chord line
	if chordCount == len(tokens) && chordCount > 0 {
		return true
	}

	// Special case: Line with significant leading spaces and only chords (chord positioning)
	// e.g., "       A" or "    D          A"
	// This is very common in chord charts where chords are positioned above lyrics
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " \t"))
	if leadingSpaces >= 3 && chordCount == len(tokens) && chordCount > 0 {
		return true
	}

	// If line has lots of internal spacing and all non-space content is chords
	// This catches lines like "    D          A" where spaces position chords
	if chordCount > 0 && float64(len(line)) >= float64(len(trimmed))*1.5 && chordCount == len(tokens) {
		return true
	}

	// Check if line contains common lyric words - if so, NOT a chord line
	if containsCommonWords(trimmed) {
		return false
	}

	// If more than 50% of tokens are chords and no common words, likely a chord line
	if chordCount > 0 && float64(chordCount)/float64(len(tokens)) >= 0.5 {
		return true
	}

	return false
}

// containsCommonWords checks if line contains common lyric words (to distinguish from chord lines)
func containsCommonWords(line string) bool {
	// Common words that appear in lyrics but not in chord lines
	commonWords := []string{
		"the", "and", "you", "your", "my", "me", "i", "a", "to", "in", "of", "is", "it",
		"for", "on", "with", "that", "this", "from", "all", "will", "can", "when", "where",
		"who", "what", "have", "has", "had", "been", "was", "were", "are", "be", "he", "she",
		"we", "they", "them", "their", "his", "her", "our", "us", "him", "there", "then",
		"but", "as", "at", "by", "an", "if", "or", "so", "up", "out", "do", "not", "like",
		"just", "now", "know", "get", "got", "make", "see", "go", "come", "take", "give",
	}

	lowerLine := strings.ToLower(line)
	tokens := strings.Fields(lowerLine)

	// Check each token
	for _, token := range tokens {
		// Remove punctuation for checking
		cleanToken := strings.Trim(token, ",.!?;:'\"")
		for _, word := range commonWords {
			if cleanToken == word {
				return true
			}
		}
	}

	return false
}
