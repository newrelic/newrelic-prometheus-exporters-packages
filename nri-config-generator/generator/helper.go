package generator

import (
	"regexp"
	"strings"
)

var (
	reTrailingCommaObject = regexp.MustCompile(`,(\s)*}`)
	reTrailingCommaList   = regexp.MustCompile(`,(\s)*]`)
)

func removeTrailingCommas(content string) string {
	content = reTrailingCommaObject.ReplaceAllString(content, `}`)
	return reTrailingCommaList.ReplaceAllString(content, `]`)
}

func compactTextInOneLine(content string) string {
	return strings.ReplaceAll(content, "\n", "")
}
