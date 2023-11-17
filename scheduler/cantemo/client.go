package cantemo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Url       string
	AuthToken string
}

func New(url, authToken string) *Client {
	return &Client{
		Url:       url,
		AuthToken: authToken,
	}
}

func (c *Client) url(parts ...string) *url.URL {
	path, _ := url.JoinPath(c.Url, parts...)

	u, _ := url.Parse(path)

	return u
}

func request[T any](client *Client, method string, path string, body io.Reader) (*T, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("AUTH-TOKEN", client.AuthToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	str, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result *T
	err = json.Unmarshal(str, &result)
	return result, err
}
