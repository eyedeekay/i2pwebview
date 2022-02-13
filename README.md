# i2pwebview

Just a simple wrapper around webview configured to use I2P. Launches an
embedded i2pd router if an I2P proxy is not available. Incomplete
but usable. Doesn't support tabbed browsing, has only a limited concept
of history, and may have unknown behaviors. Not recommended, but useful
as a demonstration of how to build a WebView based app that connects to
I2P services. To compile on Linux, install all the dependencies of:
https://github.com/webview/webview and run `make example`. Builds are
untested on all other platforms.
