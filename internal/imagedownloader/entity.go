package imagedownloader

type ImageInfo struct {
	Url   string `json:"url"`
	Error string `json:"error"`
}

type Output struct {
	DownloadedImages []ImageInfo `json:"downloaded_images"`
	SkippedImages    []ImageInfo `json:"skipped_images"`
	NotFoundImages   []ImageInfo `json:"not_found_images"`
	InvalidImages    []ImageInfo `json:"invalid_images"`
	FailedImages     []ImageInfo `json:"failed_images"`
}
