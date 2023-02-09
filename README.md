# WebP Converting Service
This is a simple web service that converts images to WebP format. It is written in Go and uses the [WebP](https://developers.google.com/speed/webp/) library from Google.

## Installation
```shell
DOCKER_BUILDKIT=1 docker build -t webp-converter . && docker run -d -p 8080:8080 webp-converter
```

## Usage
```shell
curl http://localhost:8080/?url=<source_image_url>&w=<width>&q=<quality> > image.webp
```

## Benchmark
| Benchmark                 | Iter | Time            | Memory       | Allocs        |
|---------------------------|------|-----------------|--------------|---------------|
| BenchmarkConverting1080-8 | 7    | 155808351 ns/op | 2241905 B/op | 371 allocs/op |
| BenchmarkConverting720-8  | 9    | 114244537 ns/op | 2245442 B/op | 377 allocs/op |
| BenchmarkConverting480-8  | 12   | 92763052 ns/op  | 2242224 B/op | 371 allocs/op |
| BenchmarkConverting320-8  | 14   | 81400199 ns/op  | 2247629 B/op | 374 allocs/op |
| BenchmarkConverting240-8  | 14   | 77849676 ns/op  | 2214512 B/op | 257 allocs/op |
| BenchmarkConverting64-8   | 16   | 70099018 ns/op  | 2225543 B/op | 266 allocs/op |
| BenchmarkConverting32-8   | 16   | 69596867 ns/op  | 2222243 B/op | 265 allocs/op |
