package wasmsdk

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/jssdk/transactions"
	"strconv"
	"time"
)

func HacTransferSDK() {

	/* CreateHacTransfer */
	jsGlobalRegFuncPmsString("CreateHacTransfer", func(args []string) interface{} {
		if len(args) != 4 {
			return retErr(fmt.Errorf("param num must be 4."))
		}
		acc := account.GetAccountByPrivateKeyOrPassword(args[0])
		addr, e := account.CheckReadableAddress(args[1])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid Hacash address.", args[1]))
		}
		amt, e := fields.NewAmountFromString(args[2])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid transfer amount.", args[2]))
		}
		fee, e := fields.NewAmountFromString(args[3])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid fee amount.", args[3]))
		}
		// create tx
		ctime := time.Now().Unix()
		tx := transactions.CreateOneTxOfSimpleTransfer(acc, addr, amt, fee, ctime)
		txbody, e := tx.Serialize()
		if e != nil {
			return retErr(fmt.Errorf("Create tx error: %s", e.Error()))
		}
		// success
		return retData(map[string]interface{}{
			"Amount":         amt.ToFinString(),
			"Fee":            fee.ToFinString(),
			"TxHash":         tx.Hash().ToHex(),
			"TxBody":         hex.EncodeToString(txbody),
			"PaymentAddress": acc.AddressReadable,
			"Timestamp":      strconv.FormatInt(ctime, 10),
		})
	})

}
