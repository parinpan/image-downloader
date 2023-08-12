package imagedownloader

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHTTPClient_Do(t *testing.T) {
	baseReq, _ := http.NewRequest(http.MethodGet, "https://google.com/image.jpg", nil)

	t.Run("returns error on failed response - without retry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttpClient := NewMockhttpClient(ctrl)

		client := &HTTPClient{
			BaseClient:                         mockHttpClient,
			RetryOption:                        RetryOption{},
			AcceptedImageContentTypeExtensions: nil,
		}

		// mock functions
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error"))

		resp, err := client.Do(baseReq)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("returns error on failed response - with retry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttpClient := NewMockhttpClient(ctrl)

		client := &HTTPClient{
			BaseClient: mockHttpClient,
			RetryOption: RetryOption{
				BaseDelay:   time.Duration(50) * time.Millisecond,
				MaxDelay:    time.Duration(3) * time.Second,
				MaxAttempts: 3,
			},
			AcceptedImageContentTypeExtensions: nil,
		}

		// mock functions
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error")).Times(4)

		resp, err := client.Do(baseReq)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("returns error on skipped content type header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttpClient := NewMockhttpClient(ctrl)

		client := &HTTPClient{
			BaseClient: mockHttpClient,
			RetryOption: RetryOption{
				BaseDelay:   time.Duration(50) * time.Millisecond,
				MaxDelay:    time.Duration(3) * time.Second,
				MaxAttempts: 3,
			},
			AcceptedImageContentTypeExtensions: nil,
		}

		// mock functions
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
		}, nil)

		resp, err := client.Do(baseReq)
		assert.ErrorIs(t, err, ErrSkippedContentType)
		assert.NotNil(t, resp)
	})

	t.Run("returns no error on a good response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHttpClient := NewMockhttpClient(ctrl)

		client := &HTTPClient{
			BaseClient: mockHttpClient,
			RetryOption: RetryOption{
				BaseDelay:   time.Duration(50) * time.Millisecond,
				MaxDelay:    time.Duration(3) * time.Second,
				MaxAttempts: 3,
			},
			AcceptedImageContentTypeExtensions: CommonImageContentTypeExtensions,
		}

		// mock functions
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			Header: map[string][]string{
				"Content-Type": {"image/jpeg"},
			},
		}, nil)

		resp, err := client.Do(baseReq)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
