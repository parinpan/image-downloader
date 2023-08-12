package imagedownloader

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
)

const (
	contentTypeHeaderKey = "Content-Type"
)

var (
	ErrMakeRequest   = errors.New("could not build http request")
	ErrFetchResponse = errors.New("could not fetch http response")
	ErrImageNotFound = errors.New("could not download a non-existing image")
	ErrInvalidImage  = errors.New("could not download an invalid image")
	ErrOpenImageFile = errors.New("could not create a new image file")
	ErrCopyImage     = errors.New("could not copy image into the destination path")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient   httpClient
	CreateFileFn func(name string) (*os.File, error)
	CopyFileFn   func(dst io.Writer, src io.Reader) (written int64, err error)
}

func (c *Client) DownloadImage(ctx context.Context, url string, destinationPath func(contentType string) string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(ErrMakeRequest, err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.Join(ErrFetchResponse, err)
	}

	// close body in every call made
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrImageNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return ErrInvalidImage
	}

	return c.saveImage(resp.Body, destinationPath(resp.Header.Get(contentTypeHeaderKey)))
}

func (c *Client) saveImage(body io.ReadCloser, destinationPath string) error {
	file, err := c.CreateFileFn(destinationPath)
	if err != nil {
		return errors.Join(ErrOpenImageFile, err)
	}

	// close file whenever opened
	defer file.Close()

	if _, err := c.CopyFileFn(file, body); err != nil {
		return errors.Join(ErrCopyImage, err)
	}

	return nil
}
