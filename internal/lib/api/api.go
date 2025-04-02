package api

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("invalid status code")
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")
)

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after 1st redirect
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "", ErrNotFound
		case http.StatusBadRequest:
			return "", ErrBadRequest
		case http.StatusInternalServerError:
			return "", ErrInternalServerError
		}
	}
	return resp.Header.Get("Location"), nil
}
