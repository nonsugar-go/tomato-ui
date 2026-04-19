package checkpoint

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Server string
	SID    string
	HTTP   *http.Client
}

type loginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type loginResponse struct {
	SID string `json:"sid"`
}

func InsecureClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// Login authenticates with the Check Point management server and stores the session ID in the client.
// Ref: https://sc1.checkpoint.com/documents/latest/APIs/?#web/login
func (c *Client) Login(user, password string) error {
	url := fmt.Sprintf("https://%s/web_api/login", c.Server)

	reqBody := loginRequest{
		User:     user,
		Password: password,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.SID = result.SID
	return nil
}

// Post sends a POST request to the Check Point management server with the given path and body, including the session ID in the headers.
// Ref: https://sc1.checkpoint.com/documents/latest/APIs/?#web/show-hosts
func (c *Client) Post(path string, body any) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/web_api/%s", c.Server, path)

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.SID)

	return c.HTTP.Do(req)
}
