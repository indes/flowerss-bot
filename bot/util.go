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
		regexp.MustCompile("[\t\f\r ]+").ReplaceAllString(
			strip.StripTags(regexp.MustCompile("< *br */* *>").ReplaceAllString(html.UnescapeString(desc), "\n")),
			" "),
		"\n ")

	contentDescRune := []rune(desc)
	descLen := len(contentDescRune)
	// 在拉丁字母中，len(str) == len([]rune)
	// 在这里将拉丁字母长度计 0.5（在这里乘2）
	if len(desc) == descLen {
		descLen /= 2
		limit *= 2
	}
	if descLen > limit {
		desc = string(contentDescRune[:limit]) + "…"
	}

	return desc
}
