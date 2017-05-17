package funcs

import (
	"strings"
	"unicode"
)

func Camal(s string) string {
	if strings.IndexRune(s, '_') < 0 {
		return strings.ToLower(s)
	}
	return strings.Map(camalFunctor(), s)
}

func camalFunctor() func(rune) rune {
	b := true
	return func(r rune) rune {
		if b {
			b = false
			return unicode.ToUpper(r)
		}

		if r == '_' {
			b = true
			return -1
		}

		return unicode.ToLower(r)
	}
}
