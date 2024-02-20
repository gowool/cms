package http

import (
	"html"
	"regexp"
	"strings"
)

var (
	re      = regexp.MustCompile(`<(.|\n)*?>`)
	escaper = strings.NewReplacer(`"`, "&quot;")
)

func StripTags(content string) string {
	return re.ReplaceAllString(html.UnescapeString(content), "")
}

func EscapeDoubleQuotes(content string) string {
	return escaper.Replace(content)
}
