// Webassembly project main.go
package main

import "github.com/hacash/jssdk/wasmsdk"

/**

// build WASM file
// build base58 for WASM module
GOARCH=wasm GOOS=js go build -o dist/hacash_sdk.wasm main/main.go && go run build/main.go

// tinygo
tinygo build -o dist/hacash_sdk_tiny.wasm -target wasm  main/main.go && go run build/main.go -tiny=true

// wasm zip
zip -j -9 ./dist/hacash_sdk_tiny.zip ./dist/hacash_sdk_tiny.wasm && cp ./dist/hacash_sdk_tiny.zip ./webwallet/lib/hacash_sdk_tiny.zip

// build webwallet zip
zip -q -r -v -9 hacash_web_wallet.zip webwallet -x'*.less'

// buile single file
go run build/main.go -build_single_html=true && zip -j -o hacash_web_wallet_single_page.zip hacash_wallet_single_page.html



*/

func main() {
	// hacash js SDKs

	wasmsdk.AccountSDK()
	wasmsdk.HacTransferSDK()
	wasmsdk.SatTransferSDK()
	wasmsdk.HacdTransferSDK()

	select {}
}
