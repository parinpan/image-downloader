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
		urls = append(urls, scanner.Text())

		if len(urls) == cap(urls) {
			if err := batchExecutor(urls); err != nil {
				return err
			}
			// clear urls after processed
			urls = urls[:0]
		}
	}

	// execute remaining urls
	if err := batchExecutor(urls); err != nil {
		return err
	}

	return nil
}