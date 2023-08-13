export IMAGE=$1

if [ -z "$IMAGE_STORE_PATH" ]; then
  echo "ERROR: You must specify IMAGE_STORE_PATH env when running this script"
  exit 1;
fi;

if [ "$INPUT_FIXTURE" ]; then
  docker run -v "$INPUT_FIXTURE":/fixtures/images.txt -v "$IMAGE_STORE_PATH":/downloads "$IMAGE"
else
  docker run -v "$IMAGE_STORE_PATH":/downloads "$IMAGE"
fi
