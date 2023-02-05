package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	app := iris.New()
	app.Get("/", webpConverter)
	app.Listen(":8080")
}

func webpConverter(ctx iris.Context) {
	hash := GetMD5Hash(ctx.Request().URL.RawQuery)
	url := ctx.URLParams()["url"]
	if url == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	width := ctx.URLParams()["w"]
	if width == "" {
		width = "0"
	}

	quality := ctx.URLParams()["q"]
	if quality == "" {
		quality = "80"
	}
	_ = hash
	temp, err := os.CreateTemp("", "*")
	defer os.Remove(temp.Name())
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	getImageFromUrl(url, temp)
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	f, err := os.OpenFile("./cache/"+hash, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	cmd := exec.Command(cWebpBin, "-q", quality, temp.Name(), "-o", "-", "-resize", width, "0", "-m", "6", "-mt")
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	out, err := cmd.StdoutPipe()
	defer out.Close()
	if err != nil {
		throwInternalError(ctx, err)
		return
	}
	if err := cmd.Start(); err != nil {
		throwInternalError(ctx, err)
		return
	}

	ctx.Header("Content-Type", "image/webp")
	ctx.StreamWriter(func(w io.Writer) bool {
		mWriter := io.MultiWriter(w, f)
		_, err := io.Copy(mWriter, out)
		return err != nil
	})
	cmd.Wait()
}

func getImageFromUrl(url string, writer io.Writer) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, &bytes.Buffer{})

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func throwInternalError(ctx iris.Context, err error) {
	ctx.StatusCode(iris.StatusInternalServerError)
	log.Println(err)
}
