package handlers

import (
	"regexp"
	"strings"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/pkg/chordanalysis"
)

// wrapChordsInBrackets converts plain text chords to bracketed format for analysis
// Handles formats like "D" "Am" "F#m" "Bm7" "C/E" on their own or with spacing
func wrapChordsInBrackets(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	// Chord pattern: matches standalone chords (not part of words)
	// Matches: D, Am, F#m, Bm7, C/E, Dsus4, etc.
	chordPattern := regexp.MustCompile(`\b([A-G][#b]?(?:maj|min|m|sus|aug|dim|add)?\d*(?:/[A-G][#b]?)?)\b`)

	for _, line := range lines {
		// Check if line is primarily chords (more than 50% of non-space content is chords)
		// or if line contains spaced out chords
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and section headers
		if trimmed == "" || strings.HasSuffix(trimmed, ":") {
			result = append(result, line)
			continue
		}

		// Skip lines that already have brackets (from worshipchords or other sources)
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			result = append(result, line)
			continue
		}

		// Count how many chord matches we have
		matches := chordPattern.FindAllString(trimmed, -1)

		// If line has chords and looks like a chord line (short words, mostly chords)
		// Replace chords with bracketed versions
		if len(matches) > 0 {
			// Check if this looks like a chord line vs lyrics
			// Chord lines typically have: short length, multiple spaces, no lowercase words
			words := strings.Fields(trimmed)
			hasLowercase := false
			for _, word := range words {
				if len(word) > 0 && word[0] >= 'a' && word[0] <= 'z' {
					hasLowercase = true
					break
				}
			}

			// If no lowercase and has chords, it's likely a chord line
			if !hasLowercase || len(words) <= 3 {
				// Wrap all chords in brackets
				newLine := chordPattern.ReplaceAllString(line, "[$1]")
				result = append(result, newLine)
				continue
			}
		}

		// Default: keep line as-is
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// FormatWithNashville formats content to OnSong format with chord analysis
// This ensures all sources (Ultimate Guitar, Worshipchords, Manual) use consistent formatting
func FormatWithNashville(songName, artistName, key, content string) string {
	var output strings.Builder

	// Wrap plain text chords in brackets for analysis
	content = wrapChordsInBrackets(content)

	// Clean up content - convert section headers
	sectionPattern := regexp.MustCompile(`(?mi)^\[(Intro|Verse\s*\d*|Chorus\s*\d*|Pre-Chorus|Bridge|Instrumental|Interlude|Turnaround|Outro|Tag|Ending|Solo|Break|Refrain|Coda|Hook|Vamp|Outro Chorus)\]\s*$`)
	content = sectionPattern.ReplaceAllString(content, "$1:")

	// Perform chord analysis
	analysis := chordanalysis.AnalyzeChords(content)

	// Header: Song name, artist
	output.WriteString(songName + "\n")
	output.WriteString(artistName + "\n")

	// Use detected key if available, otherwise use provided key
	keyToUse := key
	if analysis.DetectedKey != "" {
		keyToUse = analysis.DetectedKey
	}
	// Fallback if no key at all
	if keyToUse == "" {
		keyToUse = "C"
	}

	output.WriteString("Key: " + keyToUse + "\n")
	output.WriteString("Tempo: 100 BPM\n")
	output.WriteString("Time Signature: 4/4\n")
	output.WriteString("\n")

	// Add the content
	output.WriteString(content)

	return output.String()
}
