package tokenizer

import (
	"strings"
	"unicode"
	"golang.org/x/text/unicode/norm"
)

func StringNormalize(s string) string {
	// TODO: Remove stop words

	// Convert to lowercase
	s = strings.ToLower(s)

	// Normalize to remove accents (NFD form splits characters from accents)
	s = norm.NFD.String(s)

	// Remove diacritics
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Mn category is for non-spacing marks
			return -1
		}
		return r
	}, s)

	// Remove punctuation and special characters, retain spaces and alphanumerics
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	s = b.String()

	// Trim spaces
	s = strings.TrimSpace(s)

	// Replace multiple spaces with a single space
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// nGrams generates a slice of n-grams from the provided document string.
func NGrams(document string, n int) []string {
	document = StringNormalize(document)
	ngrams := []string{}
	documentLength := len(document)
	for i := 0; i <= documentLength-n; i++ {
		ngram := document[i : i+n]
		ngrams = append(ngrams, ngram)
	}
	return ngrams
}

func NGramSet(document string, n int) map[string]struct{} {
	nGrams := NGrams(document, n)
	nGramSet := make(map[string]struct{})
	for _, nGram := range nGrams {
		nGramSet[nGram] = struct{}{}
	}
	return nGramSet
}

func NGramFrequency(document string, n int) map[string]int {
	nGrams := NGrams(document, n)
	nGramFrequency := make(map[string]int)
	for _, nGram := range nGrams {
		nGramFrequency[nGram]++
	}
	return nGramFrequency
}