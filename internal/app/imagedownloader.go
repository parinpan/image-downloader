package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/oklog/ulid/v2"

	"fachr.in/image-downloader/internal/fixture"
	"fachr.in/image-downloader/internal/imagedownloader"
	"fachr.in/image-downloader/internal/util"
	imageDownloaderPkg "fachr.in/image-downloader/pkg/imagedownloader"
)

func StartImageDownloaderApp(ctx context.Context) error {
	imageDownloader := &imagedownloader.ImageDownloader{
		FixtureLoader: &fixture.Fixture{
			Path:      "fixtures/images.txt",
			BatchSize: 25,
		},
		DownloaderClient: &imageDownloaderPkg.Client{
			HTTPClient: &imageDownloaderPkg.HTTPClient{
				BaseClient: http.DefaultClient,
				RetryOption: imageDownloaderPkg.RetryOption{
					BaseDelay:   time.Duration(50) * time.Millisecond,
					MaxDelay:    time.Duration(3) * time.Second,
					MaxAttempts: 3,
				},
				AcceptedImageContentTypeExtensions: imageDownloaderPkg.CommonImageContentTypeExtensions,
			},
			CreateFileFn: os.Create,
			CopyFileFn:   io.Copy,
		},
		UlidMakerFn:                      ulid.Make,
		Workers:                          10,
		StorageRootPath:                  "downloads/",
		CommonImageContentTypeExtensions: imageDownloaderPkg.CommonImageContentTypeExtensions,
	}

	out, err := imageDownloader.DownloadAllImages(ctx)
	if err != nil {
		return err
	}

	return util.JsonStdout(out)
}
