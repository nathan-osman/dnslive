package client

import (
	"io"
	"net/http"
)

const (
	ipv4CheckUrl = "https://ipv4.icanhazip.com"
	ipv6CheckUrl = "https://ipv6.icanhazip.com"
)

func get(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) update() error {

	return nil
}
