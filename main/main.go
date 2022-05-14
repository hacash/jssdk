// Webassembly project main.go
package main

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"syscall/js"
)

/**

// build WASM file
// build base58 for WASM module

GOARCH=wasm GOOS=js go build -o dist/hacash_sdk.wasm main/main.go && go run build/main.go



*/

func foo(that js.Value, args []js.Value) interface{} {
	fmt.Println("hellow wasm")
	fmt.Println(args)
	return nil
}

func main() {
	// 将golang中foo函数注入到window.foo中
	js.Global().Set("foo", js.FuncOf(foo))
	// 将100注入到 window.value中
	js.Global().Set("value", 100)

	// hacash js sdk

	var retAcc = func(acc *account.Account) interface{} {
		v := js.ValueOf(map[string]interface{}{
			"Address":    acc.AddressReadable,
			"PublicKey":  hex.EncodeToString(acc.PublicKey),
			"PrivateKey": hex.EncodeToString(acc.PrivateKey),
		})
		return v
	}

	/* CreateNewRandomAccount */
	js.Global().Set("CreateNewRandomAccount", js.FuncOf(func(that js.Value, args []js.Value) interface{} {
		acc := account.CreateNewRandomAccount()
		return retAcc(acc)
	}))

	/* GetAccountByPrivateKeyOrPassword */
	js.Global().Set("GetAccountByPrivateKeyOrPassword", js.FuncOf(func(that js.Value, args []js.Value) interface{} {
		acc := account.GetAccountByPrivateKeyOrPassword(args[0].String())
		return retAcc(acc)
	}))

	select {}
}
