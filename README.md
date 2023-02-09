# WebP Converting Service
This is a simple web service that converts images to WebP format. It is written in Go and uses the [WebP](https://developers.google.com/speed/webp/) library from Google.

## Installation
```shell
DOCKER_BUILDKIT=1 docker build -t webp-converter . && docker run -d -p 8080:8080 webp-converter
```