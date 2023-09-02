# High-Level Design
![imagedownloader](https://github.com/parinpan/image-downloader/assets/14908455/7ef1ad2f-a8db-480b-b697-d04d301c62e8)


## Description
1. The ImageDownloaderApp initiates and triggers the ImageDownloaderService.
2. The ImageDownloaderService retrieves batched image URLs from the FixtureLoaderExecutor.
3. The ImageDownloaderService then distributes these batched image URLs among its worker pool.
4. Each worker within the ImageDownloaderService processes its assigned batch of image URLs.
5. Utilizing the ImageDownloaderClient, each ImageDownloaderWorker efficiently downloads the specified image URLs.
6. The downloaded images are stored on the local disk by the ImageDownloaderClient.
7. The ImageDownloaderApp generates a report on STDOUT, detailing downloaded images, unavailable images, skipped images, and more.

# Key Strengths of This Solution
Several underlying implementations set this solution apart:

1. The ImageDownloaderService employs a worker pool comprising 10 workers, with each worker capable of executing 25 concurrent download operations. This parallel approach accelerates image downloading.
2. Through connection pooling, the ImageDownloaderClient optimizes HTTP requests by avoiding the overhead of establishing new connections for each call. This results in significantly reduced latency.
3. The ImageDownloaderClient incorporates an HTTP retry mechanism, allowing failed calls to be retried up to three times, enhancing the solution's robustness.
4. A strategic exponential backoff strategy is applied to the retry mechanism in the ImageDownloaderClient, contributing to improved reliability in the face of connectivity challenges.
5. Utilizing HTTP timeouts, the ImageDownloaderClient prevents application hang-ups due to unexpectedly prolonged tasks, enhancing overall responsiveness.
6. Image IDs are generated using ULID, ensuring that images with the same name do not overwrite each other, thus maintaining data integrity.
7. Instead of relying solely on image extensions in URLs, the ImageDownloaderClient identifies image types based on the content type header. This versatile approach ensures accurate identification irrespective of URL structures.
8. The application is containerized using Docker, making it portable and runnable in diverse environments. A setup script is provided, simplifying the deployment process. Additionally, users can customize input fixtures and image storage paths to suit their requirements.

# How To

This section explains how to interact with the app.

## Running Tests

To run the tests, execute the following command in your terminal:
```
make test
```

## Run the app
### Prerequisite

Before proceeding, make sure you have Docker installed on your machine. If you're using an Apple Silicon machine, you can run Docker using Colima. Refer to this guide for installation and usage instructions: [Colima GitHub](https://github.com/abiosoft/colima)

### Building the app
Before running the app. Please build it first. Execute this command in the terminal:
```bash
make build
```

### Running with Default Input Fixture
To run the app with the default input fixture, follow these steps in your terminal. Make sure that the image store directory exists and is writable:
```bash
IMAGE_STORE_PATH=/absolute/path/to/store/your/image/directory make run
```

### Running with Custom Input Fixture
To run the app with a custom input fixture, follow these steps in your terminal. Ensure that the input file exists and is in the same format as the default fixture (check **fixtures/images.txt** in this repository). Confirm that the image store directory exists and is writable:
```bash
INPUT_FIXTURE=/absolute/path/to/your/input/images.txt IMAGE_STORE_PATH=/absolute/path/to/store/your/image/directory make run
```
