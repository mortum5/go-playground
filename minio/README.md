# Minio

File uploader on Fiber

## Install

1. Install [docker](https://www.docker.com/) and docker compose

## Run

1. Copy env file with `cp env_example .env`
2. Run minio instance with `docker compose up`
3. Upload file with command `curl -i -X POST -H "Content-Type: multipart/form-data" -F "fileUpload=@PATH_TO_FILE" http://localhost:5000/api/v1/upload`
4. Open `http://localhost:9001` in browser