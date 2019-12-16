package bot

import (
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
			regexp.MustCompile(`\n+`).ReplaceAllLiteralString(
				strings.ReplaceAll(
					regexp.MustCompile(`<br(| /)>`).ReplaceAllString(desc, "<br>"),
					"<br>", "\n"),
				"\n")),
		"\n")

	contentDescRune := []rune(desc)
	if len(contentDescRune) > limit {
		desc = string(contentDescRune[:limit])
	}

	return desc
}
