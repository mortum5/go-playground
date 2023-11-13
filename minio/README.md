# File uploader on Fiber

```sh
docker run \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio1 \
  -v $(PWD)/data:/data \
  -e "MINIO_ROOT_USER=AKIAIOSFODNN7EXAMPLE" \
  -e "MINIO_ROOT_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  quay.io/minio/minio server /data --console-address ":9001"
```

Import .env var

```sh
set -o allexport
source .env 
set +o allexport
```