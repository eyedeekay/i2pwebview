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
				go func() {
					if err := exec.Command(ex, args...).Run(); err != nil {
						log.Fatal(err)
					}
				}()
			}
		}
	}
}

func webView(u string) {
	ex, err := os.Executable()
	args := []string{"-w=" + strconv.Itoa(*width), "-h=" + strconv.Itoa(*height), "-p=" + *proxy, "-d=" + strconv.FormatBool(*debug), u}
	log.Println(ex, args, err)
	w := webview.New(*debug)
	defer w.Destroy()
	w.SetTitle("I2P Surf")
	w.SetSize(*width, *height, webview.HintNone)
	//	w.Init()
	w.Navigate(`data:text/html,<html><body>
	<!--<label for="back">Back</label><input type="button" id="back" name="back" value="Back"> 
	<label for="next">Next</label><input type="button" id="next" name="next" value="Next">-->
	<label for="furl">URL:</label><input width="` + strconv.FormatFloat(float64(*width)*.8, 'f', -1, 64) + `px" type="text" id="furl" name="furl">
	<label for="go">Go</label><input type="button" id="go" name="go" onclick="go()" value="Go">
	<script type="text/javascript">function go(){
	  if (document.getElementById('furl').value != "") {
	    document.getElementById('htmlview').src = document.getElementById('furl').value
	  }
	}
	</script>
	<iframe id="htmlview" width="` + strconv.Itoa(*width) + `px" height="` + strconv.Itoa(*height) + `px" src="` + u + `"><iframe>
	</body></html>`)
	w.Run()
}
