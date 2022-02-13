package main

import (
	"flag"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eyedeekay/go-i2pcontrol"
	i2pd "github.com/eyedeekay/go-i2pd/goi2pd"
	"github.com/webview/webview"
)

var (
	width         = flag.Int("w", 1280, "width of the window to create the webview in")
	height        = flag.Int("h", 1024, "height of the window to create the webview in")
	proxy         = flag.String("p", "127.0.0.1:4444", "I2P HTTP proxy to use when browsing")
	routerconsole = flag.String("r", "127.0.0.1:7657", "I2P router console to use when browsing")
	debug         = flag.Bool("d", true, "Debug mode")
)

func proxyTest() bool {
	_, err := net.Listen("tcp", *proxy)
	if err != nil {
		// proxy is up
		return true
	}
	return false
}

func proxyLoop() bool {
	for {
		if proxyTest() {
			return true
		}
	}
}

func main() {
	flag.Parse()
	if !proxyTest() {
		closer := i2pd.InitI2PSAM()
		defer closer()
		i2pd.StartI2P()
		defer i2pd.StopI2P()
		*routerconsole = "127.0.0.1:7070"
	}
	log.Println("Waiting for proxy to come up")
	if proxyLoop() {
		log.Println("Proxy is up")
	}
	os.Setenv("http_proxy", "http://"+*proxy)
	os.Setenv("https_proxy", "http://"+*proxy)
	os.Setenv("ftp_proxy", "http://"+*proxy)
	os.Setenv("all_proxy", "http://"+*proxy)
	os.Setenv("HTTP_PROXY", "http://"+*proxy)
	os.Setenv("HTTPS_PROXY", "http://"+*proxy)
	os.Setenv("FTP_PROXY", "http://"+*proxy)
	os.Setenv("ALL_PROXY", "http://"+*proxy)
	os.Setenv("no_proxy", "http://"+*routerconsole)
	os.Setenv("NO_PROXY", "http://"+*routerconsole)
	log.Println("Proxy set for:", *proxy)
	log.Println("Router console set for:", *routerconsole)
	if len(flag.Args()) == 0 {
		log.Println("Opening default page")
		webView("http://" + *routerconsole)
	} else if len(flag.Args()) == 1 {
		log.Println("Opening single page")
		webView(flag.Args()[0])
	} else {
		log.Println("Opening multiple pages")
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

func goToAddress() string {
	return `
	let history = [];
	let index = 0;
	function goTo(){
		let url = document.getElementById("furl").value;
		if (history.length == 0){
			history = [url];
		}
		goURL(url);
	}
	function goURL(url){
		console.log("Going to:", url);
		if (url != "") {
			if ( !url.startsWith("http://") ){
				console.log("didn't start with http", url);
				if ( !url.startsWith("https://") ) {
					console.log("didn't start with https", url);
					let httpUrl = "http://".concat(url);
					url = httpUrl;
					console.log("Force-set URL", url, httpUrl);
				}
			}
			if (!history.includes(url))
				history = [url, ...history];
			document.getElementById('htmlview').src = url;
			console.log("History:", history, index)
		}
		index = history.length;
	}
	function goBack(){
		if (index === 0 ) {
			index == history.length;
		}
		index--;
		let url = history[index];
		goURL(url);
	}
	function goForward(){
		if (index === history.length ) {
			index = 0;
		}
		index++;
		let url = history[index];
		goURL(url);
	}
	/*function updatenetworkStatus(){
		let status = document.getElementById("networkstatus");
		netStatus().then((x) => {
			status.innerText = x;
			console.log("Network status: ", status.innerText);
		});
	}
	function updateStatus(){
		let status = document.getElementById("routerStatus");
		softStatus().then((x) => {
			status.innerText = xs;
			console.log("Router status: ", status.innerText);
		});
	}
	function updateTunnels(){
		let tunnels = document.getElementById("participatingCount");
		participating().then((x) => {
			tunnels.innerText = x;
			console.log("Tunnel count: ", tunnels.innerText);
		});
	}*/
	//setInterval(updatenetworkStatus, 1000);
	//setInterval(updateStatus, 1000);
	//setInterval(updateTunnels, 1000);
	window.onload = function(e) {
		var input = document.getElementById("furl");
		input.addEventListener("keyup", function(event) {
			if (event.keyCode === 13) {
				console.log("Going to a URL");
				event.preventDefault();
				document.getElementById("go").click();
			}
		});
	}
`
}

func webView(u string) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	log.Println("Opening page:", u)
	ex, err := os.Executable()
	args := []string{"-w=" + strconv.Itoa(*width), "-h=" + strconv.Itoa(*height), "-p=" + *proxy, "-d=" + strconv.FormatBool(*debug), u}
	log.Println(ex, args, err)

	w := webview.New(*debug)
	defer w.Destroy()

	w.SetTitle("I2P Surf")
	w.SetSize(*width, *height, webview.HintNone)

	w.Bind("netStatus", i2pcontrol.NetStatus)
	w.Bind("softStatus", i2pcontrol.Status)
	w.Bind("participating", i2pcontrol.ParticipatingTunnels)
	w.Bind("restart", i2pcontrol.RestartGraceful)
	w.Bind("shutdown", i2pcontrol.ShutdownGraceful)

	w.Navigate(`data:text/html,<html><body>
	<style>
	.mini {
		font-size: 9px;
		min-width: 10%;
	}
	label {
		border: 1px solid black !important;
		border-radius: 5px;
		padding: 5px;
		margin: 2px;
	}
	span {
		width: 80%;
	}
	.padded {
		padding: 9px;
	}
	.window {
		width: 100%;
		height: 100%;
		overflow: hidden;
		min-width: 100%;
		min-height: 100%;
		width: 100vw;
		height: 100vh;
		overflow: hidden;
		min-width: 100vw;
		min-height: 100vh;
	}
	.tabbar {
		width: 100%;
		overflow: hidden;
		min-width: 100%;
		width: 100vw;
		overflow: hidden;
		min-width: 100vw;
		min-height: 1em;
		border: 1px solid black !important;
		border-radius: 5px;
		padding: 6px;
		margin: 2px;
	}
	.navbar {
		min-width: 30%;
	}
	</style>
	<script type="text/javascript">
	` + goToAddress() + `
	</script>
	<div class="tabbar">
	<label class="mini padded" for="back">Back
		<input type="button" id="back" name="back" value="Back" onclick="goBack()">
	</label>
	<label class="mini padded" for="next">Next
		<input type="button" id="next" name="next" value="Next" onclick="goForward()">
	</label>
	<label class="navbar" for="furl">URL:
		<input width="` + strconv.FormatFloat(float64(*width)*.8, 'f', -1, 64) + `px" type="text" id="furl" name="furl" value="` + u + `">
	</label>
	<!--<label for="go">Go</label>-->
	<input type="button" id="go" name="go" onclick="goTo()" value="Go">
	<label class="mini" for="netstatus">Network Status:
		<span class="mini" id="networkStatus">Network Status</span>
	</label>
	<label class="mini" for="status">Status:
		<span class="mini" id="routerStatus">Router Status</span>
	</label>
	<label class="mini" for="participating">Participating Tunnels:
		<span class="mini" id="participatingCount">Participating Tunnels</span>
	</label>
	</div>
	<div class="tabbar">}}
		<span class="tabgroup"></span>
	</div>
	<div class="window">` + iframe(*width, *height, u) + `</div>
	</body></html>`)
	//
	w.Run()
}

func iframe(width, height int, url string) string {
	return `<iframe id="htmlview" class="window" width="100vw" height="100vh" src="` + url + `"><iframe>`
	//return `<iframe id="htmlview" class="window" width="` + strconv.Itoa(width) + `px" height="` + strconv.Itoa(height) + `px" src="` + u + `"><iframe>``
}
