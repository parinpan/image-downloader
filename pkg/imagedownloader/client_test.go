package imagedownloader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClient_DownloadImage(t *testing.T) {
	ctx := context.Background()

	destinationPath := func(contentType string) string {
		return "absolute/path/to/image"
	}

	acceptedContentTypes := map[string]string{
		"image/jpeg": ".jpg",
	}

	t.Run("returns error on an invalid request", func(t *testing.T) {
		client := &Client{}
		err := client.DownloadImage(ctx, ":) !some_invalid_url! :)", nil)
		assert.ErrorIs(t, err, ErrMakeRequest)
	})

	t.Run("returns error on an invalid response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)
		client := Client{HTTPClient: mockHttp}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error"))

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", nil)
		assert.ErrorIs(t, err, ErrFetchResponse)
	})

	t.Run("returns error when content-type skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"something/else"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", destinationPath)
		assert.ErrorIs(t, err, ErrSkippedContentType)
	})

	t.Run("returns error on 404 http status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", nil)
		assert.ErrorIs(t, err, ErrImageNotFound)
	})

	t.Run("returns error on any non 200 http status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", nil)
		assert.ErrorIs(t, err, ErrInvalidImage)
	})

	t.Run("returns error when couldn't create file to store image", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
			CreateFileFn: func(name string) (*os.File, error) {
				return nil, errors.New("error")
			},
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", destinationPath)
		assert.ErrorIs(t, err, ErrOpenImageFile)
	})

	t.Run("returns error when couldn't copy image", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
			CreateFileFn: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			CopyFileFn: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, errors.New("error")
			},
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", destinationPath)
		assert.ErrorIs(t, err, ErrCopyImage)
	})

	t.Run("returns no error when everything is ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttp := NewMockhttpClient(ctrl)

		client := Client{
			HTTPClient:                         mockHttp,
			AcceptedImageContentTypeExtensions: acceptedContentTypes,
			CreateFileFn: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			CopyFileFn: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
		}

		// mock http response
		mockHttp.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
			Body: io.NopCloser(bytes.NewBuffer(nil)),
		}, nil)

		err := client.DownloadImage(ctx, "https://fachr.in/static/image/fachrin-memoji.jpg", destinationPath)
		assert.NoError(t, err)
	})
}
