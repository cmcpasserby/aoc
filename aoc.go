package aoc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	sessionCookie string
	baseUrl       *url.URL
	httpClient    *http.Client
}

func New(sessionCookie string, httpClient *http.Client) *Client {
	return &Client{
		sessionCookie: sessionCookie,
		baseUrl:       &url.URL{Scheme: "https", Host: "adventofcode.com"},
		httpClient:    httpClient,
	}
}

func (c *Client) DownloadInput(year, day int) (io.Reader, error) {
	rel := &url.URL{Path: fmt.Sprintf("%d/day/%d/input", year, day)}
	u := c.baseUrl.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name: "session",
		Value:  c.sessionCookie,
	}
	req.AddCookie(cookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
