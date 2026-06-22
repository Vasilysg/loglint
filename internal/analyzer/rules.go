package analyzer

import (
	"strings"
	"unicode"
)

func startsWithLower(s string) bool {
	if s == "" {
		return true
	}

	first := []rune(s)[0]
	return unicode.IsLower(first)
}

func isEnglishText(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}

func hasOnlyAllowedChars(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			continue
		}

		return false
	}

	return true
}

func hasSensitiveWords(s string, patterns []string) bool {
	s = strings.ToLower(s)

	for _, pattern := range patterns {
		pattern = strings.ToLower(pattern)

		if strings.Contains(s, pattern) {
			return true
		}
	}

	return false
}
