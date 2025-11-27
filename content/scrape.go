package content

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	lg "github.com/charmbracelet/lipgloss"
	readability "github.com/go-shiori/go-readability"
)

var (
	h1Style = lg.NewStyle().
		Foreground(lg.Color("147")).
		Bold(true).
		Underline(true)

	h2Style = lg.NewStyle().
		Foreground(lg.Color("147")).
		Bold(true)

	h3Style = lg.NewStyle().
		Foreground(lg.Color("111")).
		Bold(true)

	headerStyle = lg.NewStyle().
			Foreground(lg.Color("111")).
			Bold(true)

	codeBlockStyle = lg.NewStyle().
			Foreground(lg.Color("229")).
			Background(lg.Color("235")).
			Padding(1, 2)

	codeInlineStyle = lg.NewStyle().
			Foreground(lg.Color("223")).
			Background(lg.Color("236"))
)

func GetArticleContent(articleURL string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", articleURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// sneaky
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// parse with readability
	parsedURL, _ := url.Parse(articleURL)
	article, err := readability.FromReader(strings.NewReader(string(htmlBytes)), parsedURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse article: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var output strings.Builder

	// process elements
	doc.Find("body").Contents().Each(func(i int, s *goquery.Selection) {
		renderNode(&output, s)
	})

	return output.String(), nil
}

func renderNode(output *strings.Builder, s *goquery.Selection) {
	nodeName := goquery.NodeName(s)

	switch nodeName {
	case "p":
		text := strings.TrimSpace(s.Text())
		if text != "" {
			output.WriteString(text)
			output.WriteString("\n\n")
		}

	case "h1", "h2", "h3", "h4", "h5", "h6":
		text := strings.TrimSpace(s.Text())
		if text != "" {
			output.WriteString("\n")
			switch nodeName {
			case "h1":
				output.WriteString(h1Style.Render(text))
			case "h2":
				output.WriteString(h2Style.Render(text))
			case "h3":
				output.WriteString(h3Style.Render(text))
			default:
				output.WriteString(headerStyle.Render(text))
			}
			output.WriteString("\n\n")
		}

	case "pre":
		code := s.Text()
		if code != "" {
			output.WriteString("\n")
			lines := strings.SplitSeq(code, "\n")
			for line := range lines {
				output.WriteString(codeBlockStyle.Render(line))
				output.WriteString("\n")
			}
			output.WriteString("\n")
		}

	case "code":
		if s.Parent().Length() > 0 && goquery.NodeName(s.Parent()) != "pre" {
			output.WriteString(codeInlineStyle.Render(s.Text()))
		}

	case "ul", "ol":
		s.Children().Each(func(i int, li *goquery.Selection) {
			if goquery.NodeName(li) == "li" {
				prefix := "•"
				if nodeName == "ol" {
					prefix = fmt.Sprintf("%d.", i+1)
				}
				text := strings.TrimSpace(li.Text())
				if text != "" {
					fmt.Fprintf(output, "  %s %s\n", prefix, text)
				}
			}
		})
		output.WriteString("\n")

	case "img":
		alt, hasAlt := s.Attr("alt")
		if hasAlt && alt != "" {
			fmt.Fprintf(output, "\n[Image: %s]\n\n", alt)
		}

	case "br":
		output.WriteString("\n")

	case "blockquote":
		text := strings.TrimSpace(s.Text())
		if text != "" {
			lines := strings.SplitSeq(text, "\n")
			for line := range lines {
				if strings.TrimSpace(line) != "" {
					output.WriteString("│ " + line + "\n")
				}
			}
			output.WriteString("\n")
		}

	case "#text":
		// only render text nodes that are direct children of body
		// (cuz text inside other elements is handled by those elements)
		parent := s.Parent()
		if parent.Length() > 0 {
			parentName := goquery.NodeName(parent)
			if parentName == "body" || parentName == "div" {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					output.WriteString(text)
					output.WriteString("\n")
				}
			}
		}

	default:
		// just recurse into children for errything else
		s.Contents().Each(func(i int, child *goquery.Selection) {
			renderNode(output, child)
		})
	}
}
