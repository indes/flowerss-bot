package preview

import (
	"html"
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

func TrimDescription(desc string, limit int) string {
	if limit == 0 {
		return ""
	}
	desc = strings.Trim(
		strip.StripTags(
			regexp.MustCompile("\n+").ReplaceAllLiteralString(
				strings.ReplaceAll(
					regexp.MustCompile(`<br(| /)>`).ReplaceAllString(
						html.UnescapeString(desc), "<br>",
					),
					"<br>", "\n",
				),
				"\n",
			),
		),
		"\n",
	)

	contentDescRune := []rune(desc)
	if len(contentDescRune) > limit {
		desc = string(contentDescRune[:limit])
	}

	return desc
}
