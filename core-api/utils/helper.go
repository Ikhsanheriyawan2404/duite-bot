package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func Slugify(s string) string {
    s = strings.ToLower(s)
    var b strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            b.WriteRune(r)
        } else if unicode.IsSpace(r) || r == '-' {
            b.WriteRune('-')
        }
        // abaikan karakter lain
    }
    return strings.Trim(b.String(), "-")
}

func ContainsNominal(s string) bool {
	re := regexp.MustCompile(`\d+`)
	return re.MatchString(s)
}
