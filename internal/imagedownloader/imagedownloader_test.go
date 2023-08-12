package imagedownloader

import (
	"context"
	"errors"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"fachr.in/image-downloader/internal/fixture"
	"fachr.in/image-downloader/pkg/imagedownloader"
)

func TestImageDownloader_DownloadAllImages(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error when could not load fixture", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFixture := NewMockfixtureLoader(ctrl)

		imageDownloader := &ImageDownloader{
			FixtureLoader: mockFixture,
			Workers:       3,
		}

		// mock functions
		mockFixture.EXPECT().LoadExecute(ctx, gomock.Any()).Return(errors.New("error"))

		out, err := imageDownloader.DownloadAllImages(ctx)
		assert.Error(t, err)
		assert.Equal(t, (*Output)(nil), out)
	})

	t.Run("returns no error when fixture is loaded and executed successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDownloaderClient := NewMockdownloaderClient(ctrl)

		imageDownloader := &ImageDownloader{
			FixtureLoader: &fixture.Fixture{
				Path:      "./testdata/images.txt",
				BatchSize: 20,
			},
			DownloaderClient:                 mockDownloaderClient,
			UlidMakerFn:                      ulid.Make,
			Workers:                          3,
			StorageRootPath:                  "/some/storage/path",
			CommonImageContentTypeExtensions: imagedownloader.CommonImageContentTypeExtensions,
		}

		// mock functions
		mockDownloaderClient.EXPECT().DownloadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(imagedownloader.ErrSkippedContentType)
		mockDownloaderClient.EXPECT().DownloadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(imagedownloader.ErrImageNotFound)
		mockDownloaderClient.EXPECT().DownloadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(imagedownloader.ErrInvalidImage)
		mockDownloaderClient.EXPECT().DownloadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		out, err := imageDownloader.DownloadAllImages(ctx)
		assert.NoError(t, err)
		assert.Len(t, out.DownloadedImages, 1)
		assert.Len(t, out.FailedImages, 1)
		assert.Len(t, out.InvalidImages, 1)
		assert.Len(t, out.NotFoundImages, 1)
	})
}

func TestImageDownloader_destinationPath(t *testing.T) {
	t.Run("returns path based on url and content type", func(t *testing.T) {
		imageDownloader := &ImageDownloader{
			UlidMakerFn: func() (id ulid.ULID) {
				return ulid.MustNew(0, nil)
			},
			StorageRootPath:                  "/downloads",
			CommonImageContentTypeExtensions: imagedownloader.CommonImageContentTypeExtensions,
		}

		assert.Equal(t, "/downloads/00000000000000000000000000_a.jpg", imageDownloader.destinationPath("https://a.com/a")("image/jpeg"))
		assert.Equal(t, "/downloads/00000000000000000000000000_a.jpg", imageDownloader.destinationPath("https://a.com/a.jpeg")("image/jpeg"))
		assert.Equal(t, "/downloads/00000000000000000000000000_%2F.jpg", imageDownloader.destinationPath("https://a.com/")("image/jpeg"))
		assert.Equal(t, "/downloads/00000000000000000000000000_.jpg", imageDownloader.destinationPath("https://a.com")("image/jpeg"))
	})
}
