package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func copyFile(srcFile, destFile string) (int64, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()
	return io.Copy(file2, file1)
}

// build single html
func buildSingleHTML() {
	pwd, _ := os.Getwd()
	fmt.Println("buildSingleHTML pwd work dir:", pwd)

	var readFileMust = func(fp string) []byte {
		hfp := path.Join(pwd, fp)
		fbts, e := os.ReadFile(hfp)
		if e != nil {
			err := fmt.Errorf("%s cannot find: %s", hfp, e.Error())
			fmt.Println(err.Error())
			panic(err)
		}
		return fbts
	}


	// read html file
	htmlFileName := "./webwallet/hacash_wallet.html"
	html := readFileMust(htmlFileName)
	// split & insert js code
	htmlstr := string(html)
	htmlstr = strings.Replace(htmlstr, `<link rel=\"stylesheet\" href='./lib/index.css'>`, "", 1)
	htmlstr = strings.Replace(htmlstr, `<script src="./lib/wasm_exec.js" type="text/javascript"></script>`, "", 1)
	htmlstr = strings.Replace(htmlstr, `<script src="./lib/hacash_sdk.js" type="text/javascript"></script>`, "", 1)
	htmlstr = strings.Replace(htmlstr, `<script src="./lib/index.js" type="text/javascript"></script>`, "", 1)

	// merge
	htmlstrAry := strings.Split(htmlstr, "<!--code-insert-->")
	resultHTML := make([]string, 7)
	resultHTML[0] = htmlstrAry[0]

	// index.css
	css1 := readFileMust("./webwallet/lib/index.css")
	resultHTML[1] = "<style>" + string(css1) + "</style>"

	// <!--code-insert-->
	resultHTML[2] = htmlstrAry[1]

	// wasm_exec.js
	js1 := readFileMust("./webwallet/lib/wasm_exec.js")
	resultHTML[3] = "<script>" + string(js1) + "</script>"

	// hacash_sdk.js
	js2 := readFileMust("./webwallet/lib/hacash_sdk.js")
	resultHTML[4] = "<script>" + string(js2) + "</script>"

	// index.js
	js3 := readFileMust("./webwallet/lib/index.js")
	resultHTML[5] = "<script>" + string(js3) + "</script>"

	// <!--code-insert-->
	resultHTML[6] = htmlstrAry[2]

	// save single file
	page1 := path.Join(pwd, "./hacash_wallet_single_page.html")
	os.WriteFile(page1, []byte(strings.Join(resultHTML, "\n")), os.ModePerm)
	fmt.Printf("create file: %s\n", page1)

}

// build wasm
func buildWASM(usetinyfix string) {

	pwd, _ := os.Getwd()
	fmt.Println("buildWASM pwd work dir:", pwd)

	// ad wasm file
	wasmFileName := "./dist/hacash_sdk" + usetinyfix + ".wasm"
	wfp := path.Join(pwd, wasmFileName)
	wasm, e := os.ReadFile(wfp)
	if e != nil {
		err := fmt.Errorf("%s cannot find: %s", wasmFileName, e.Error())
		fmt.Println(err.Error())
		panic(err)
	}

	// wasm to base64
	wasmString := base64.StdEncoding.EncodeToString(wasm)
	wasmJsContent := `var hacash_sdk_wasm_code_base64 = "` + wasmString + `";
		function base64ToBuffer(b) {
			const str = window.atob(b);
			const buffer = new Uint8Array(str.length);
			for (let i=0; i < str.length; i++) {
				buffer[i] = str.charCodeAt(i);
			}
			return buffer;
		}
		var hacash_sdk_wasm_code = base64ToBuffer(hacash_sdk_wasm_code_base64);
	`

	sdk1 := path.Join(pwd, "./webwallet/lib/hacash_sdk.js")
	os.WriteFile(sdk1, []byte(wasmJsContent), os.ModePerm)
	fmt.Printf("create %s => %s\n", wfp, sdk1)

	// copy ./dist/wasm_exec.js
	exec1 := path.Join(pwd, "./dist/wasm_exec"+usetinyfix+".js")
	exec2 := path.Join(pwd, "./webwallet/lib/wasm_exec.js")
	copyFile(exec1, exec2)
	fmt.Printf("copy %s => %s\n", exec1, exec2)

	fmt.Println("build successfully!")

}



func main() {

	tiny := flag.Bool("tiny", false, "use tinygo")
	build_single_html := flag.Bool("build_single_html", false, "build all html css js file to single html")
	flag.Parse()

	if *build_single_html {
		// build wallet static html
		buildSingleHTML()
	}else{
		// build wasm
		if *tiny {
			fmt.Println("[Use tinygo compiling]")
			buildWASM("_tiny")
		} else {
			buildWASM("")
		}

	}
}
