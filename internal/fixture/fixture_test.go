package fixture

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixture_LoadExecute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error on an opening file failure", func(t *testing.T) {
		fixture := &Fixture{
			Path:      "open/non/existing/file/txt.txt",
			BatchSize: 1,
		}

		err := fixture.LoadExecute(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("returns error on failed batch execution", func(t *testing.T) {
		fixture := &Fixture{
			Path:      "./testdata/images.txt",
			BatchSize: 1,
		}

		batchExecutor := func(urls []string) error {
			return errors.New("error")
		}

		err := fixture.LoadExecute(ctx, batchExecutor)
		assert.Error(t, err)
	})

	t.Run("returns error on failed batch execution - with remaining batch", func(t *testing.T) {
		fixture := &Fixture{
			Path:      "./testdata/images.txt",
			BatchSize: 2,
		}

		batchExecutor := func(urls []string) error {
			if len(urls) != fixture.BatchSize {
				return errors.New("error")
			}
			return nil
		}

		err := fixture.LoadExecute(ctx, batchExecutor)
		assert.Error(t, err)
	})

	t.Run("returns no error on succeeded batch execution", func(t *testing.T) {
		fixture := &Fixture{
			Path:      "./testdata/images.txt",
			BatchSize: 1,
		}

		var collectedUrls []string

		expectedUrls := []string{
			"https://a.com/a.jpg",
			"https://b.com/c.png",
			"https://c.com/c.gif",
		}

		batchExecutor := func(urls []string) error {
			collectedUrls = append(collectedUrls, urls...)
			return nil
		}

		err := fixture.LoadExecute(ctx, batchExecutor)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedUrls, collectedUrls)
	})

	t.Run("returns no error on succeeded batch execution - with remaining batch", func(t *testing.T) {
		fixture := &Fixture{
			Path:      "./testdata/images.txt",
			BatchSize: 2,
		}

		var collectedUrls []string

		expectedUrls := []string{
			"https://a.com/a.jpg",
			"https://b.com/c.png",
			"https://c.com/c.gif",
		}

		batchExecutor := func(urls []string) error {
			collectedUrls = append(collectedUrls, urls...)
			return nil
		}

		err := fixture.LoadExecute(ctx, batchExecutor)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedUrls, collectedUrls)
	})
}
