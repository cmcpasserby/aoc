package aoc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const baseUrl = "adventofcode.com"

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

func NewClient(sessionId string) (*Client, error) {
	u := &url.URL{Scheme: "https", Host: baseUrl}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies := []*http.Cookie{
		{
			Name:  "session",
			Value: sessionId,
		},
	}
	jar.SetCookies(u, cookies)

	return &Client{
		baseUrl:    u,
		httpClient: &http.Client{Jar: jar},
	}, nil
}

func (c *Client) GetInputs(year, day int) (*Puzzle, error) {
	path := fmt.Sprintf("%d/day/%d/input", year, day)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	inputs, err := c.do(req)

	puzzle := &Puzzle{
		Year:  year,
		Day:   day,
		Input: inputs,
	}

	return puzzle, err
}

func (c *Client) Submit(answer *Answer) (error, string) {
	if err := answer.Validate(); err != nil {
		return err, ""
	}

	path := fmt.Sprintf("%d/day/%d/answer", answer.Year, answer.Day)

	form := &url.Values{}
	form.Set("level", strconv.Itoa(int(answer.Part)))
	form.Set("answer", answer.Answer)


	req, err := c.newRequest("POST", path, strings.NewReader(form.Encode()))
	if err != nil {
		return err, ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := c.do(req)
	if err != nil {
		return err, ""
	}

	message, err := ioutil.ReadAll(body)
	if err != nil {
		return err, ""
	}
	msg := string(message)

	result := ""

	if strings.Contains(msg, "That's the right answer") {
		result = fmt.Sprintf("the answer %s is correct", answer.Answer)
	} else if strings.Contains(msg, "Did you already complete it") {
		result = "Already completed!"
	} else if strings.Contains(msg, "That's not the right answer") {
		result = fmt.Sprintf("The answer %s is INCORRECT", answer.Answer)
	} else if strings.Contains(msg, "You gave an answer too recently") {
		// TODO parse wait information
	}

	return err, result
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseUrl.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (io.Reader, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		message, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code: %d, message:\n%q", resp.StatusCode, string(message))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(body)

	return buf, nil
}
