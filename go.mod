module github.com/eyedeekay/i2pwebview

go 1.17

require (
	github.com/eyedeekay/go-i2pcontrol v0.0.0-20201227222421-6e9f31a29a91
	github.com/eyedeekay/go-i2pd v0.0.0-20220213070306-9807541b2dfc
	github.com/webview/webview v0.0.0-20220212050339-62c361a22430
)

require github.com/ybbus/jsonrpc/v2 v2.1.6 // indirect

replace github.com/eyedeekay/go-i2pd v0.0.0-20220213070306-9807541b2dfc => ./go-i2pd
