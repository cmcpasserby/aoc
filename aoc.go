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

func (c *Client) newRequest(path string) (*http.Request, error) {
	rel := &url.URL{Path: path}
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

	return req, nil
}

func (c *Client) do(req *http.Request) (io.Reader, error) {
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

func (c *Client) GetInput(year, day int) (io.Reader, error) {
	req, err := c.newRequest(fmt.Sprintf("%d/day/%d/input", year, day))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) GetQuestion(year, day int) (io.Reader, error) {
	req, err := c.newRequest(fmt.Sprintf("%d/day/%d", year, day))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
