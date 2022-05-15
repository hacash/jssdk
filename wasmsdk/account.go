package wasmsdk

import (
	"github.com/hacash/core/account"
	"syscall/js"
)


func AccountSDK() {

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

}
