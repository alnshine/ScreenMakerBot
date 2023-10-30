package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gosimple/slug"
)

type Alphabet struct {
}

func main() {
	startTime := time.Now()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	file, err := os.Open("site.txt")
	if err != nil {
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := "https://" + scanner.Text() + "/"
		var filename string
		if len(os.Args) == 3 {
			filename = os.Args[2]
		} else {
			filename = slug.Make(url) + ".png"
		}

		var imageBuf []byte

		if err := chromedp.Run(
			ctx,
			ScreenshotTasks(url, &imageBuf),
		); err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
			log.Fatal(err)
		}
	}
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)

	fmt.Printf("Время выполнения: %v\n", elapsedTime)
}

func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().Do(ctx)
			return err
		}),
	}
}
