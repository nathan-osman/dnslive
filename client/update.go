package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ipv4CheckUrl = "https://ipv4.icanhazip.com"
	ipv6CheckUrl = "https://ipv6.icanhazip.com"
)

type updateParams struct {
	Name string `json:"name"`
	Ipv4 string `json:"ipv4,omitempty"`
	Ipv6 string `json:"ipv6,omitempty"`
}

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

// Only one of the IPv4 or IPv6 lookups needs to succeed; only if both fail
// does the function call fail. In this case, the returned error is from the
// first call (IPv4).

func (c *Client) update() error {
	var (
		params = updateParams{
			Name: c.name,
		}
		encounteredErr error
	)
	if v, err := get(ipv4CheckUrl); err != nil {
		encounteredErr = err
	} else {
		params.Ipv4 = strings.TrimSpace(v)
	}
	if v, err := get(ipv6CheckUrl); err != nil {
		if encounteredErr != nil {
			return err
		}
	} else {
		params.Ipv6 = strings.TrimSpace(v)
	}
	u := url.URL{
		Scheme: "https",
		Host:   c.serverAddr,
		Path:   "/update",
	}
	b, err := json.Marshal(&params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return err
}
