package content

import (
	"html"
	"regexp"
	"strings"
)

// StripHTML remove HTML tags and decodes HTML entities from string
func StripHTML(s string) string {
	s = regexp.MustCompile(`(?is)<script.*?</script>`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`(?is)<style.*?</style>`).ReplaceAllString(s, "")

	// decode entities like &amp, &lt, etc
	s = html.UnescapeString(s)

	re := regexp.MustCompile(`(?s)<.*?>`)
	s = re.ReplaceAllString(s, "")

	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}
