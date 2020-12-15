// Package remote contains methods for interacting with remote providers.
package remote

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
	ipld "github.com/ipfs/go-ipld-format"
)

// Remote is used to interact with external services.
type Remote struct {
	url    string
	client *http.Client
}

// NewRemote returns a remote using the given url.
func NewRemote(url string) *Remote {
	return &Remote{
		url:    url,
		client: http.DefaultClient,
	}
}

// Upload converts the dags starting at the given roots into CAR format and uploads it.
func (r *Remote) Upload(ctx context.Context, dag ipld.DAGService, roots ...cid.Cid) error {
	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)

	fileWriter, err := bodyWriter.CreateFormFile("file", "")
	if err != nil {
		return err
	}

	if err := car.WriteCar(ctx, dag, roots, fileWriter); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", r.url, "api/v0/dag/import")
	contentType := bodyWriter.FormDataContentType()

	if err := bodyWriter.Close(); err != nil {
		return err
	}

	resp, err := http.Post(url, contentType, &body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("upload failed status=%s reply=%s", resp.Status, string(reply))
	}

	return nil
}
