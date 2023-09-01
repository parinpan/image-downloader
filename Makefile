.PHONY: build

test:
	go test ./...

build:
	docker build . -f build/Dockerfile -t fachrin/image-downloader:latest

run:
	INPUT_FIXTURE=$(INPUT_FIXTURE) IMAGE_STORE_PATH=$(IMAGE_STORE_PATH) ./scripts/docker_run.sh fachrin/image-downloader:latest
