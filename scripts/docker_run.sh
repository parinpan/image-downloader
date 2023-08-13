export IMAGE=$1

if [ "$INPUT_FIXTURE" ]; then
  docker run -v "$INPUT_FIXTURE":/fixtures/images.txt -v "$IMAGE_STORE_PATH":/downloads "$IMAGE"
else
  docker run -v "$IMAGE_STORE_PATH":/downloads "$IMAGE"
fi
