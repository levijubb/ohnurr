package reader

import (
	"io"
	"net/http"
)

func GetArticleContent(url string) (string, error) {
	r, err := http.Get(url)

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
