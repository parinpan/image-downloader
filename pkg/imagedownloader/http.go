package imagedownloader

import (
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var (
	ErrSkippedContentType = errors.New("skip image due to not listed in accepted content type")
)

type RetryOption struct {
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	MaxAttempts int
}

type HTTPClient struct {
	BaseClient                         httpClient
	RetryOption                        RetryOption
	AcceptedImageContentTypeExtensions map[string]string
}

func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := h.do(req, 0)
	if err != nil {
		return resp, err
	}

	contentType := resp.Header.Get(contentTypeHeaderKey)

	if _, ok := h.AcceptedImageContentTypeExtensions[contentType]; !ok {
		return resp, errors.Join(ErrSkippedContentType, err)
	}

	return resp, nil
}

func (h *HTTPClient) do(req *http.Request, retryCount int) (*http.Response, error) {
	delay := rand.Float64() * math.Pow(2.0, float64(retryCount)) * float64(h.RetryOption.BaseDelay)
	delayDuration := time.Duration(delay)

	if retryCount != 0 {
		time.Sleep(delayDuration)
	}

	resp, err := h.BaseClient.Do(req)
	retryable := err != nil && delayDuration <= h.RetryOption.MaxDelay && retryCount+1 <= h.RetryOption.MaxAttempts

	if retryable {
		return h.do(req, retryCount+1)
	}

	return resp, err
}
