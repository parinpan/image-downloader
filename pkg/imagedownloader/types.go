package imagedownloader

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	CommonImageContentTypeExtensions = map[string]string{
		"image/jpeg":         ".jpg",
		"image/png":          ".png",
		"image/gif":          ".gif",
		"image/bmp":          ".bmp",
		"image/webp":         ".webp",
		"image/svg+xml":      ".svg",
		"image/x-icon":       ".ico",
		"image/tiff":         ".tiff",
		"image/vnd.radiance": ".hdr",
		"image/jp2":          ".jp2",
	}
)
