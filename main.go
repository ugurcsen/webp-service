package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	_url "net/url"
	"os"
	"strconv"
)

var fileSizeLimit int64
var port string

var sizes = [12]int{32, 64, 240, 320, 480, 640, 720, 1080, 1440, 2160, 2880, 4320}

func main() {
	// Getting env
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fileSizeLimitStr := os.Getenv("FILE_SIZE_LIMIT")
	if fileSizeLimitStr == "" {
		fileSizeLimit = 5 << 20 // 5MB
	} else {
		var err error
		fileSizeLimit, err = strconv.ParseInt(fileSizeLimitStr, 10, 64)
		if err != nil {
			panic(err)
		}
		fileSizeLimit = fileSizeLimit << 10
	}
	log.Println("file size limit: ", fileSizeLimit>>10, "KB")

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app := iris.Default()
	app.Use(crs)
	app.Use(iris.Compression)
	app.Get("/", webpConverter)
	app.Listen(":" + port)
}

func webpConverter(ctx iris.Context) {
	// Validating
	url, err := _url.Parse(ctx.URLParam("url"))
	if err != nil {
		ctx.ContentType("text/plain")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	width, err := strconv.Atoi(ctx.URLParam("w"))
	if err != nil {
		width = 0
	} else {
		flag := false
		for _, size := range sizes {
			if width == size {
				flag = true
				break
			}
		}
		if !flag {
			ctx.ContentType("text/plain")
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("width must be one of these: 32, 64, 240, 320, 480, 640, 720, 1080, 1440, 2160, 2880, 4320")
			return
		}
	}

	quality, err := strconv.Atoi(ctx.URLParam("q"))
	if err != nil {
		quality = 80
	} else {
		if quality < 0 || quality > 100 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("quality must be a number between 0 and 100")
			return
		}
	}

	// Downloading
	temp, err := os.CreateTemp("", "*")
	defer os.Remove(temp.Name())
	defer temp.Close()
	if err != nil {
		throwInternalError(ctx, err)
		return
	}

	err = getImageFromUrl(url, temp, fileSizeLimit)
	if err != nil {
		ctx.ContentType("text/plain")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	// Converting
	out, err := convertWebp(ctx.Request().Context(), temp.Name(), quality, width)
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	defer out.Close()

	// Writing
	ctx.ContentType("image/webp")
	_, err = io.Copy(ctx.ResponseWriter(), out)
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
}
