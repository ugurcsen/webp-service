package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
)

func throwInternalError(ctx iris.Context, err error) {
	ctx.StatusCode(iris.StatusInternalServerError)
	log.Println(err)
}

func getImageFromUrl(url *url.URL, writer io.Writer, limit int64) error {
	res, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.ContentLength > limit {
		return fmt.Errorf("file size is too big. limit: %dKB, got: %dKB", limit>>10, res.ContentLength>>10)
	}

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func convertWebp(ctx context.Context, tempFile string, quality int, width int) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, cWebpBin, "-q", strconv.Itoa(quality), tempFile, "-o", "-", "-resize", strconv.Itoa(width), "0", "-m", "6", "-mt")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	return out, nil
}
