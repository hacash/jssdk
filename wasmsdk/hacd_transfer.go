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

func HacdTransferSDK() {

	/* CreateHacTransfer */
	jsGlobalRegFuncPmsString("CreateHacdTransfer", func(args []string) interface{} {
		if len(args) != 5 {
			return retErr(fmt.Errorf("param num must be 5."))
		}
		acc := account.GetAccountByPrivateKeyOrPassword(args[0])
		addr, e := account.CheckReadableAddress(args[1])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid Hacash address.", args[1]))
		}
		diamonds := fields.NewEmptyDiamondListMaxLen200()
		e = diamonds.ParseHACDlistBySplitCommaFromString(args[2])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid hacd names, ERROR: %s", args[2], e.Error()))
		}
		acc_fee := account.GetAccountByPrivateKeyOrPassword(args[3])
		fee, e := fields.NewAmountFromString(args[4])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid fee amount.", args[3]))
		}
		// create tx
		ctime := time.Now().Unix()
		tx, e := transactions.CreateOneTxOfOutfeeQuantityHACDTransfer(acc, addr, args[2], acc_fee, fee, ctime)
		if e != nil {
			return retErr(fmt.Errorf("create tx error: %s", e.Error()))
		}
		txbody, e := tx.Serialize()
		if e != nil {
			return retErr(fmt.Errorf("Create tx error: %s", e.Error()))
		}
		// success
		return retData(map[string]interface{}{
			"Diamonds":       diamonds.SerializeHACDlistToCommaSplitString(),
			"DiamondCount":   int(diamonds.Count),
			"Fee":            fee.ToFinString(),
			"TxHash":         tx.Hash().ToHex(),
			"TxBody":         hex.EncodeToString(txbody),
			"FeeAddress": 	  acc_fee.AddressReadable,
			"PaymentAddress": acc.AddressReadable,
			"CollectionAddress": args[1],
			"Timestamp":      strconv.FormatInt(ctime, 10),
		})
	})

}
