package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"fachr.in/image-downloader/internal/app"
)

func main() {
	cliApp := cli.App{
		Name: "start",
		Action: func(ctx *cli.Context) error {
			return app.StartImageDownloaderApp(ctx.Context)
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
