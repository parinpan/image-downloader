package imagedownloader

type Output struct {
	DownloadedImages []string `json:"downloaded_images"`
	SkippedImages    []string `json:"skipped_images"`
	NotFoundImages   []string `json:"not_found_images"`
	InvalidImages    []string `json:"invalid_images"`
	FailedImages     []string `json:"failed_images"`
}
