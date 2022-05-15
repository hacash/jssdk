package wasmsdk

import (
	"encoding/hex"
	"github.com/hacash/core/account"
	"syscall/js"
)

func retAcc(acc *account.Account) interface{} {
	v := js.ValueOf(map[string]interface{}{
		"Address":    acc.AddressReadable,
		"PublicKey":  hex.EncodeToString(acc.PublicKey),
		"PrivateKey": hex.EncodeToString(acc.PrivateKey),
	})
	return v
}

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
