// Webassembly project main.go
package main

import "github.com/hacash/jssdk/wasmsdk"

/**

// build WASM file
// build base58 for WASM module
GOARCH=wasm GOOS=js go build -o dist/hacash_sdk.wasm main/main.go && go run build/main.go

// tinygo
tinygo build -o dist/hacash_sdk_tiny.wasm -target wasm  main/main.go && go run build/main.go -tiny=true


// build webwallet zip
cd webwallet && zip -q -r -v -9 ../hacash_web_wallet.zip * && cd ..



*/

func main() {
	// hacash js SDKs

	wasmsdk.AccountSDK()

	select {}
}
