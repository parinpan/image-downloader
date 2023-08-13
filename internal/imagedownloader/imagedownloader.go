package imagedownloader

import (
	"context"
	"fmt"
	uri "net/url"
	"path"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"

	"fachr.in/image-downloader/pkg/imagedownloader"
	"fachr.in/image-downloader/pkg/logger"
)

type downloaderClient interface {
	DownloadImage(ctx context.Context, url string, destinationPath func(contentType string) string) error
}

type fixtureLoader interface {
	LoadExecute(ctx context.Context, batchExecutor func(urls []string) error) error
}

type ImageDownloader struct {
	FixtureLoader                    fixtureLoader
	DownloaderClient                 downloaderClient
	UlidMakerFn                      func() (id ulid.ULID)
	Workers                          int
	StorageRootPath                  string
	CommonImageContentTypeExtensions map[string]string
}

func (i *ImageDownloader) DownloadAllImages(ctx context.Context) (*Output, error) {
	jobs := make(chan []string)
	result := make(chan Output)

	// spawn workers
	for worker := 0; worker < i.Workers; worker++ {
		go i.worker(ctx, worker, jobs, result)
	}

	enqueueJobs := func(urls []string) error {
		jobs <- urls
		return nil
	}

	if err := i.FixtureLoader.LoadExecute(ctx, enqueueJobs); err != nil {
		return nil, err
	}

	close(jobs)
	out := <-result

	return &out, nil
}

func (i *ImageDownloader) downloadImages(ctx context.Context, urls []string) Output {
	var out Output
	var wg sync.WaitGroup

	for _, url := range urls {
		if _, err := uri.ParseRequestURI(url); err != nil {
			out.InvalidImages = append(out.InvalidImages, ImageInfo{
				Url:   url,
				Error: "image url is invalid",
			})
			continue
		}

		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			err := i.DownloaderClient.DownloadImage(ctx, url, i.destinationPath(url))

			imageInfo := ImageInfo{
				Url: url,
			}

			if err != nil {
				imageInfo.Error = err.Error()
				logger.Errorf("could not download image, imageInfo: %v", imageInfo)
			}

			switch err {
			case nil:
				logger.Infof("image downloaded: %v", imageInfo)
				out.DownloadedImages = append(out.DownloadedImages, imageInfo)
			case imagedownloader.ErrImageNotFound:
				out.NotFoundImages = append(out.NotFoundImages, imageInfo)
			case imagedownloader.ErrSkippedContentType:
				out.SkippedImages = append(out.SkippedImages, imageInfo)
			default:
				out.FailedImages = append(out.FailedImages, imageInfo)
			}
		}(url)
	}

	// wait until all images downloaded
	wg.Wait()

	return out
}

func (i *ImageDownloader) worker(ctx context.Context, id int, jobs chan []string, result chan Output) {
	for urls := range jobs {
		if len(urls) > 0 {
			logger.Infof("worker: %d - downloading images from %v", id, urls)
			result <- i.downloadImages(ctx, urls)
		}
	}
}

func (i *ImageDownloader) destinationPath(url string) func(string) string {
	u, _ := uri.Parse(url)
	id := i.UlidMakerFn().String()

	imageName := path.Base(u.Path)
	imageNameWithoutExt := strings.TrimSuffix(imageName, path.Ext(imageName))
	fileName := fmt.Sprintf("%s_%s", imageNameWithoutExt, strings.ToLower(id))

	return func(contentType string) string {
		ext := i.CommonImageContentTypeExtensions[contentType]
		return fmt.Sprintf("%s/%s%s", i.StorageRootPath, uri.PathEscape(fileName), ext)
	}
}
