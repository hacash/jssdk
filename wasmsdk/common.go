package wasmsdk

import (
	"encoding/hex"
	"github.com/hacash/core/account"
	"syscall/js"
)

func jsGlobalRegFunc(funcName string, funcBody func(that js.Value, args []js.Value) interface{}) {
	js.Global().Set(funcName, js.FuncOf(funcBody))
}

func jsGlobalRegFuncPmsString(funcName string, funcBody func(params []string) interface{}) {
	jsGlobalRegFunc(funcName, func(that js.Value, args []js.Value) interface{} {
		var params = make([]string, len(args))
		for i, v := range args {
			params[i] = v.String()
		}
		return funcBody(params)
	})
}

/* ----------------- */

func retErr(err error) interface {} {
	return js.ValueOf(map[string]interface{}{
		"Error": err.Error(),
	})
}

func retSingleData(key string, value interface {}) interface {} {
	return js.ValueOf(map[string]interface{}{
		key: value,
	})
}

func retData(data map[string]interface{}) interface {} {
	return js.ValueOf(data)
}


func retAcc(acc *account.Account) interface{} {
	v := js.ValueOf(map[string]interface{}{
		"Address":    acc.AddressReadable,
		"PublicKey":  hex.EncodeToString(acc.PublicKey),
		"PrivateKey": hex.EncodeToString(acc.PrivateKey),
	})
	return v
}
