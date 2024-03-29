package aoc

import (
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

func (c *Client) newRequest(path string) (*http.Request, error) {
	u := c.baseUrl.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "github.com/cmcpasserby/aoc by chris@rightsomegoodgames.ca")

	cookie := &http.Cookie{
		Name:  "session",
		Value: c.sessionCookie,
	}
	req.AddCookie(cookie)

	return req, nil
}

func (c *Client) do(req *http.Request) ([]byte, error) {
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

	return b, nil
}

func (c *Client) GetInput(year, day int) ([]byte, error) {
	req, err := c.newRequest(fmt.Sprintf("%d/day/%d/input", year, day))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
