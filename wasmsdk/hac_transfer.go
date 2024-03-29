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
		if len(args) != 5 {
			return retErr(fmt.Errorf("param num must be 5."))
		}
		chain_id, e := strconv.ParseInt(args[0], 10, 64)
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid chain id.", args[0]))
		}
		acc := account.GetAccountByPrivateKeyOrPassword(args[1])
		addr, e := account.CheckReadableAddress(args[2])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid Hacash address.", args[2]))
		}
		amt, e := fields.NewAmountFromString(args[3])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid transfer amount.", args[3]))
		}
		fee, e := fields.NewAmountFromString(args[4])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid fee amount.", args[4]))
		}
		// create tx
		ctime := time.Now().Unix()
		tx := transactions.CreateOneTxOfSimpleTransfer(uint64(chain_id), acc, addr, amt, fee, ctime)
		txbody, e := tx.Serialize()
		if e != nil {
			return retErr(fmt.Errorf("Create tx error: %s", e.Error()))
		}
		// success
		return retData(map[string]interface{}{
			"Amount":            amt.ToFinString(),
			"Fee":               fee.ToFinString(),
			"TxHash":            tx.Hash().ToHex(),
			"TxBody":            hex.EncodeToString(txbody),
			"PaymentAddress":    acc.AddressReadable,
			"CollectionAddress": args[2],
			"Timestamp":         strconv.FormatInt(ctime, 10),
		})
	})

}
