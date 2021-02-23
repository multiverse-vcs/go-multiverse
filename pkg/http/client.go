package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Client is an HTTP client.
type Client struct {
	*http.Client
}

// NewClient returns a new client.
func NewClient() *Client {
	return &Client{
		&http.Client{},
	}
}

// Fetch returns the repository at the given remote.
func (c *Client) Fetch(remote string) (*object.Repository, error) {
	url := fmt.Sprintf("http://%s/%s", BindAddr, remote)

	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("fetch request failed")
	}

	var repo object.Repository
	if err := json.NewDecoder(res.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

// Push updates the given branch at the remote with data.
func (c *Client) Push(remote string, branch string, data []byte) error {
	url := fmt.Sprintf("http://%s/%s/%s", BindAddr, remote, branch)

	res, err := c.Post(url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("push request failed")
	}

	return nil
}
