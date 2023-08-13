package fixture

import (
	"bufio"
	"context"
	"os"
)

type Fixture struct {
	Path      string
	BatchSize int
}

func (f *Fixture) LoadExecute(_ context.Context, batchExecutor func(urls []string) error) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	urls := make([]string, 0, f.BatchSize)

	for scanner.Scan() {
		url := scanner.Text()

		if url == "" {
			continue
		}

		// append only non-empty url to batch urls
		urls = append(urls, url)

		if len(urls) == cap(urls) {
			// batchExecutor might run in a go routine; so copy url values to make it thread safe
			var safeUrls = make([]string, len(urls))
			copy(safeUrls, urls)

			if err := batchExecutor(safeUrls); err != nil {
				return err
			}

			// clear urls after processed
			urls = urls[:0]
		}
	}

	if len(urls) == 0 {
		return nil
	}

	// execute remaining urls
	if err := batchExecutor(urls); err != nil {
		return err
	}

	return nil
}
