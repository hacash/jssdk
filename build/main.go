package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
)

func copyFile(srcFile, destFile string) (int64, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()
	return io.Copy(file2, file1)
}

func main() {

	pwd, _ := os.Getwd()
	fmt.Println("pwd work dir:", pwd)

	// ad wasm file
	wasmFileName := "./dist/hacash_sdk.wasm"
	wfp := path.Join(pwd, wasmFileName)
	wasm, e := os.ReadFile(wfp)
	if e != nil {
		err := fmt.Errorf("%s cannot find: %s", wasmFileName, e.Error())
		fmt.Println(err.Error())
		panic(err)
	}

	// wasm to base64
	wasmString := base64.StdEncoding.EncodeToString(wasm)
	wasmJsContent := `var hacashwasmcode = "` + wasmString + `";
		function base64ToBuffer(b) { 
			const str = window.atob(b); 
			const buffer = new Uint8Array(str.length); 
			for (let i=0; i < str.length; i++) { 
				buffer[i] = str.charCodeAt(i); 
			} 
			return buffer; 
		}
		var go = new Go();
		var instance = new WebAssembly.Instance(new WebAssembly.Module(base64ToBuffer(hacashwasmcode)), go.importObject); 
		go.run(instance);
		hacash_wallet_main();
	`

	file1 := path.Join(pwd, "./webwallet/lib/hacash_sdk.js")
	os.WriteFile(file1, []byte(wasmJsContent), os.ModePerm)
	fmt.Printf("create => %s\n", file1)

	// copy ./dist/wasm_exec.js
	file2 := path.Join(pwd, "./webwallet/lib/wasm_exec.js")
	copyFile(path.Join(pwd, "./dist/wasm_exec.js"), file2)
	fmt.Printf("copy => %s\n", file2)

	fmt.Println("build successfully!")

}
