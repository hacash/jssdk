// test
// document.getElementById("init").className = "ok";

var $init = document.getElementById("init")
, $init_progress = $init.getElementsByClassName("progress")[0]
, $init_bar = $init_progress.getElementsByClassName("bar")[0]
;

// loading progress
window.hacash_sdk_loading_progress = function (percent) {
    $init_progress.classList.add("show")
    $init_bar.style.width = percent + "%";
}

// check WebAssembly
var WebAssemblyIsOK = false
if(!WebAssembly || !WebAssembly.Instance || !WebAssembly.Module) {
    $init.innerHTML = "<div style='padding: 40px 20px; background: brown; color: azure; font-size: 18px; font-weight: bold;'>NOT support WebAssembly! Your browser version is too old, please upgrade your browser to the latest version.</div>"
}else{
    WebAssemblyIsOK = true
    function buildWasm(hacash_sdk_wasm_code) {
        // console.log(hacash_sdk_wasm_code)
        var go = new Go();
        WebAssembly.instantiate(hacash_sdk_wasm_code, go.importObject).then(function({module, instance}){
            go.run(instance);
            hacash_wallet_main(instance.exports);
        });
    }
    if(window.hacash_sdk_wasm_code_callback) {
        hacash_sdk_wasm_code_callback(buildWasm)
    }else if(window.hacash_sdk_wasm_code) {
        buildWasm(hacash_sdk_wasm_code)
    }
}


// start
function hacash_wallet_main(wasm){
    if(!WebAssemblyIsOK){
        return
    }

    // test
    var acc = CreateNewRandomAccount()
    console.log(acc, acc.Address)

    // var acc = GetAccountByPrivateKeyOrPassword("123456")
    // console.log(acc, acc.Address)


    // init
    init_create_account()

    // loading ok
    $init.className = "ok";

}

// create account
function init_create_account() {
    var $acc = document.getElementById("acc")
    , $ipt1 = $acc.getElementsByClassName("ipt1")[0]
    , $btn1 = $acc.getElementsByClassName("btn1")[0]
    , $btn2 = $acc.getElementsByClassName("btn2")[0]

    var in1chg = function(){
        var v = $ipt1.value
        , len = v.length
        , hav = v.replace(/[A-Za-z0-9\~\!\@\#\$\%\^\&\*\_\+\-\=\,\.\:\;]+/ig, "")
        if(len >= 6 && len <= 255 && hav == "") {
            $btn2.classList.remove("ban")
        }else{
            $btn2.classList.add("ban")
        }
    }
    $ipt1.onkeyup = in1chg
    $ipt1.onchange = in1chg
    // result
    var result = function(acc, usepass) {
        var $res = $acc.getElementsByClassName("result")[0]
        , $mark = $res.getElementsByClassName("mark")[0]
        , $b1 = $res.getElementsByClassName("addr")[0]
        , $b2 = $res.getElementsByClassName("pubkey")[0]
        , $b3 = $res.getElementsByClassName("prikey")[0]
        ;
        $mark.innerText = '( by '+ (usepass?('password '+usepass):'randomly') +' )'
        $b1.innerText = acc.Address
        $b2.innerText = acc.PublicKey
        $b3.innerText = acc.PrivateKey
        $res.classList.add("show")
    }
    // random create
    $btn1.onclick = function(){
        $ipt1.value = ""
        $btn2.classList.add("ban")
        $btn1.classList.add("ban")
        var acc = CreateNewRandomAccount()
        result(acc)
        $btn1.classList.remove("ban")
    }
    // password create
    $btn2.onclick = function(){
        if($btn2.classList.contains("ban")){
            return
        }
        var v = $ipt1.value
        , rm = v.replace(/[A-Za-z0-9\~\!\@\#\$\%\^\&\*\_\+\-\=\,\.\:\;]+/ig, "")
        if(v.length < 6 || v.length > 255){
            return alert("wrong password length")
        }
        if(rm != "") {
            return alert("unsupported symbol")
        }
        $btn1.classList.add("ban")
        var acc = GetAccountByPrivateKeyOrPassword(v)
        result(acc, v)
        $btn1.classList.remove("ban")
    }



}






