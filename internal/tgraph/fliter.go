package tgraph

import (
	"github.com/indes/flowerss-bot/internal/config"
	"regexp"
)

func match(pattern, link string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(link)
}

func Verify(link string) bool {
	if len(config.TelegraphIncludes) > 0 {
		for _, pattern := range config.TelegraphIncludes {
			if match(pattern, link) {
				return true
			}
		}
		return false
	} else if len(config.TelegraphExcludes) > 0 {
		for _, pattern := range config.TelegraphExcludes {
			if match(pattern, link) {
				return false
			}
		}
		return true
	}
	return true
}
