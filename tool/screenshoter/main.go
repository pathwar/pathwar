package main

// https://github.com/chromedp/examples/blob/master/screenshot/main.go
// https://github.com/chromedp/examples/blob/master/remote/main.go

// pull latest version of headless-shell
// $ docker pull chromedp/headless-shell:latest

// run chrome headless docker instance
// $ docker run -d -p 9222:9222 --rm --name headless-shell chromedp/headless-shell

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type chromeDebuggerJSONElement struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}
type chromeDebuggerJSON []chromeDebuggerJSONElement

func main() {
	var host, port, screenshotURL, imgQualityArg string
	flag.StringVar(&host, "host", "localhost", "Chrome headless instance ip")
	flag.StringVar(&port, "port", "9222", "Chrome headless instance debug port")
	flag.StringVar(&screenshotURL, "screenshot-url", "https://www.google.com", "URL of the page to screenshot")
	flag.StringVar(&imgQualityArg, "img-quality", "90", "PNG Quality (default : 90)")
	flag.Parse()

	// Retrieve ws url
	ipaddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get("http://" + ipaddr.IP.String() + ":" + port + "/json")
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var jsonElements chromeDebuggerJSON
	err = json.Unmarshal(jsonResponse, &jsonElements)

	if err != nil || len(jsonElements) < 1 {
		log.Fatal("Couldn't retrieve websocket headless instance URL from JSON")
	}

	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), jsonElements[0].WebSocketDebuggerURL)
	defer cancelActxt()

	ctxt, cancelCtxt := chromedp.NewContext(actxt)
	defer cancelCtxt()

	var buf []byte

	imgQuality, err := strconv.ParseInt(imgQualityArg, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if err := chromedp.Run(ctxt, fullScreenshot(screenshotURL, imgQuality, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
