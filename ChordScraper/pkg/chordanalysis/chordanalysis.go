package chordanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Chord represents a musical chord
type Chord struct {
	Root     string
	Quality  string // major, minor, seventh, etc.
	Original string
}

// Key represents a musical key
type Key struct {
	Root  string
	Scale []string // chords in the key
}

// ChordProgression represents analyzed chord progression
type ChordProgression struct {
	Chord          string
	NashvilleNum   string
	Count          int
}

// Analysis contains the full chord analysis
type Analysis struct {
	DetectedKey       string
	Confidence        float64
	AllChords         []Chord
	UniqueChords      []string
	ChordProgressions []ChordProgression
}

// Major keys and their chord scales
var majorKeys = map[string][]string{
	"C":  {"C", "Dm", "Em", "F", "G", "Am", "Bdim"},
	"C#": {"C#", "D#m", "E#m", "F#", "G#", "A#m", "B#dim"},
	"Db": {"Db", "Ebm", "Fm", "Gb", "Ab", "Bbm", "Cdim"},
	"D":  {"D", "Em", "F#m", "G", "A", "Bm", "C#dim"},
	"D#": {"D#", "E#m", "F##m", "G#", "A#", "B#m", "C##dim"},
	"Eb": {"Eb", "Fm", "Gm", "Ab", "Bb", "Cm", "Ddim"},
	"E":  {"E", "F#m", "G#m", "A", "B", "C#m", "D#dim"},
	"F":  {"F", "Gm", "Am", "Bb", "C", "Dm", "Edim"},
	"F#": {"F#", "G#m", "A#m", "B", "C#", "D#m", "E#dim"},
	"Gb": {"Gb", "Abm", "Bbm", "Cb", "Db", "Ebm", "Fdim"},
	"G":  {"G", "Am", "Bm", "C", "D", "Em", "F#dim"},
	"G#": {"G#", "A#m", "B#m", "C#", "D#", "E#m", "F##dim"},
	"Ab": {"Ab", "Bbm", "Cm", "Db", "Eb", "Fm", "Gdim"},
	"A":  {"A", "Bm", "C#m", "D", "E", "F#m", "G#dim"},
	"A#": {"A#", "B#m", "C##m", "D#", "E#", "F##m", "G##dim"},
	"Bb": {"Bb", "Cm", "Dm", "Eb", "F", "Gm", "Adim"},
	"B":  {"B", "C#m", "D#m", "E", "F#", "G#m", "A#dim"},
}

// Nashville numbering for major scale
var nashvilleNumerals = []string{"I", "ii", "iii", "IV", "V", "vi", "vii°"}

// ParseChord extracts chord information
func ParseChord(chordStr string) Chord {
	// Remove brackets if present
	chordStr = strings.Trim(chordStr, "[]")
	chordStr = strings.TrimSpace(chordStr)

	// Pattern: Root (A-G with optional # or b) + Quality (m, maj7, 7, sus4, etc.)
	re := regexp.MustCompile(`^([A-G][#b]?)(.*?)$`)
	matches := re.FindStringSubmatch(chordStr)

	if len(matches) < 3 {
		return Chord{Original: chordStr}
	}

	root := matches[1]
	quality := matches[2]

	// Determine quality type
	qualityType := "major"
	if strings.Contains(quality, "m") && !strings.Contains(quality, "maj") {
		qualityType = "minor"
	} else if strings.Contains(quality, "dim") {
		qualityType = "diminished"
	} else if strings.Contains(quality, "aug") {
		qualityType = "augmented"
	}

	return Chord{
		Root:     root,
		Quality:  qualityType,
		Original: chordStr,
	}
}

// ExtractChords finds all chords in content
func ExtractChords(content string) []Chord {
	// Pattern for chords in brackets [Chord]
	re := regexp.MustCompile(`\[([A-G][#b]?(?:maj|min|m|sus|aug|dim|add)?\d*(?:/[A-G][#b]?)?)\]`)
	matches := re.FindAllStringSubmatch(content, -1)

	var chords []Chord
	for _, match := range matches {
		if len(match) > 1 {
			chord := ParseChord(match[1])
			chords = append(chords, chord)
		}
	}

	return chords
}

// DetectKey analyzes chords and determines the most likely key
func DetectKey(chords []Chord) (string, float64) {
	if len(chords) == 0 {
		return "", 0.0
	}

	// Count chord occurrences
	chordCounts := make(map[string]int)
	for _, chord := range chords {
		baseChord := normalizeChordForKey(chord)
		chordCounts[baseChord]++
	}

	// Score each possible key
	keyScores := make(map[string]float64)

	for keyRoot, scale := range majorKeys {
		score := 0.0
		matchCount := 0

		for chordName, count := range chordCounts {
			for i, scaleChord := range scale {
				if normalizeChord(chordName) == normalizeChord(scaleChord) {
					// Weight by position in scale (I, IV, V are more important)
					weight := 1.0
					if i == 0 { // I
						weight = 3.0
					} else if i == 3 || i == 4 { // IV or V
						weight = 2.5
					} else if i == 5 { // vi (relative minor)
						weight = 2.0
					}

					score += float64(count) * weight
					matchCount++
					break
				}
			}
		}

		// Calculate confidence based on match percentage
		if matchCount > 0 {
			confidence := float64(matchCount) / float64(len(chordCounts))
			keyScores[keyRoot] = score * confidence
		}
	}

	// Find the key with highest score
	bestKey := ""
	bestScore := 0.0
	for key, score := range keyScores {
		if score > bestScore {
			bestScore = score
			bestKey = key
		}
	}

	// Calculate confidence percentage
	totalScore := 0.0
	for _, score := range keyScores {
		totalScore += score
	}

	confidence := 0.0
	if totalScore > 0 {
		confidence = (bestScore / totalScore) * 100
	}

	return bestKey, confidence
}

// GetNashvilleNumber converts a chord to Nashville number based on key
func GetNashvilleNumber(chord Chord, key string) string {
	scale, exists := majorKeys[key]
	if !exists {
		return "?"
	}

	normalizedChord := normalizeChordForKey(chord)

	for i, scaleChord := range scale {
		if normalizeChord(normalizedChord) == normalizeChord(scaleChord) {
			numeral := nashvilleNumerals[i]

			// Add modifiers
			if strings.Contains(chord.Original, "7") && !strings.Contains(chord.Original, "maj7") {
				numeral += "7"
			} else if strings.Contains(chord.Original, "maj7") {
				numeral += "M7"
			} else if strings.Contains(chord.Original, "sus") {
				numeral += "sus"
			}

			return numeral
		}
	}

	return "?"
}

// normalizeChordForKey simplifies chord to root + major/minor
func normalizeChordForKey(chord Chord) string {
	if chord.Quality == "minor" {
		return chord.Root + "m"
	}
	return chord.Root
}

// normalizeChord removes variations for comparison
func normalizeChord(chordStr string) string {
	// Remove everything after the root and m
	re := regexp.MustCompile(`^([A-G][#b]?m?)`)
	matches := re.FindStringSubmatch(chordStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return chordStr
}

// AnalyzeChords performs full chord analysis
func AnalyzeChords(content string) Analysis {
	chords := ExtractChords(content)

	// Get unique chords and count occurrences
	uniqueMap := make(map[string]int)
	for _, chord := range chords {
		uniqueMap[chord.Original]++
	}

	// Create unique list
	var uniqueChords []string
	for chord := range uniqueMap {
		uniqueChords = append(uniqueChords, chord)
	}
	sort.Strings(uniqueChords)

	// Detect key
	detectedKey, confidence := DetectKey(chords)

	// Build chord progressions with Nashville numbers
	var progressions []ChordProgression
	for _, chordName := range uniqueChords {
		chord := ParseChord(chordName)
		nashville := GetNashvilleNumber(chord, detectedKey)

		progressions = append(progressions, ChordProgression{
			Chord:        chordName,
			NashvilleNum: nashville,
			Count:        uniqueMap[chordName],
		})
	}

	// Sort by count (most common first)
	sort.Slice(progressions, func(i, j int) bool {
		return progressions[i].Count > progressions[j].Count
	})

	return Analysis{
		DetectedKey:       detectedKey,
		Confidence:        confidence,
		AllChords:         chords,
		UniqueChords:      uniqueChords,
		ChordProgressions: progressions,
	}
}

// FormatChordProgressionHeader creates a formatted header for chord progressions
func FormatChordProgressionHeader(analysis Analysis) string {
	if len(analysis.ChordProgressions) == 0 {
		return ""
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("Detected Key: %s (%.0f%% confidence)",
		analysis.DetectedKey, analysis.Confidence))
	lines = append(lines, "")
	lines = append(lines, "Chord Progressions:")

	for _, prog := range analysis.ChordProgressions {
		lines = append(lines, fmt.Sprintf("  [%s] = %s (used %dx)",
			prog.Chord, prog.NashvilleNum, prog.Count))
	}

	lines = append(lines, "")

	return strings.Join(lines, "\n")
}
