package bot

import (
	"html"
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

func trimDescription(desc string, limit int) string {
	if limit == 0 {
		return ""
	}
	desc = strings.Trim(
		strip.StripTags(
			regexp.MustCompile("\n+").ReplaceAllLiteralString(
				strings.ReplaceAll(
					regexp.MustCompile(`<br(| /)>`).ReplaceAllString(
						html.UnescapeString(desc), "<br>"),
					"<br>", "\n"),
				"\n")),
		"\n")

	contentDescRune := []rune(desc)
	descLen := len(contentDescRune)
	// in latin alphabets, len(str) == len([]rune)
	if len(desc) == descLen {
		descLen /= 2
		limit *= 2
	}
	if descLen > limit {
		desc = string(contentDescRune[:limit])
	}

	return desc
}
