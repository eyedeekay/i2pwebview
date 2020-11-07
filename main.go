package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

import "github.com/webview/webview"

var (
	width  = flag.Int("w", 800, "width of the window to create the webview in")
	height = flag.Int("h", 600, "height of the window to create the webview in")
	proxy  = flag.String("p", "127.0.0.1:4444", "I2P HTTP proxy to use when browsing")
	debug  = flag.Bool("d", false, "Debug mode")
)

func main() {
	flag.Parse()
	os.Setenv("http_proxy", "http://"+*proxy)
	os.Setenv("https_proxy", "http://"+*proxy)
	os.Setenv("ftp_proxy", "http://"+*proxy)
	os.Setenv("all_proxy", "http://"+*proxy)
	os.Setenv("HTTP_PROXY", "http://"+*proxy)
	os.Setenv("HTTPS_PROXY", "http://"+*proxy)
	os.Setenv("FTP_PROXY", "http://"+*proxy)
	os.Setenv("ALL_PROXY", "http://"+*proxy)
	if len(flag.Args()) == 1 {
		webView(flag.Args()[0])
	} else {
		for _, u := range flag.Args() {
			if _, err := url.Parse(u); err == nil {
				//        coWebView()
				ex, err := os.Executable()
				args := []string{"-w=" + strconv.Itoa(*width), "-h=" + strconv.Itoa(*height), "-p=" + *proxy, "-d=" + strconv.FormatBool(*debug), u}
				log.Println(ex, args, err)
				if err := exec.Command(ex, args...).Start(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func webView(u string) {
	w := webview.New(*debug)
	defer w.Destroy()
	w.SetTitle("I2P WebView")
	w.SetSize(*width, *height, webview.HintNone)
	w.Navigate(u)
	w.Run()
}

