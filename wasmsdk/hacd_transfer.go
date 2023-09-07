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
		if len(args) != 6 {
			return retErr(fmt.Errorf("param num must be 6."))
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
		diamonds := fields.NewEmptyDiamondListMaxLen200()
		e = diamonds.ParseHACDlistBySplitCommaFromString(args[3])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid hacd names, ERROR: %s", args[3], e.Error()))
		}
		acc_fee := account.GetAccountByPrivateKeyOrPassword(args[4])
		fee, e := fields.NewAmountFromString(args[5])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid fee amount.", args[5]))
		}
		// create tx
		ctime := time.Now().Unix()
		tx, e := transactions.CreateOneTxOfOutfeeQuantityHACDTransfer(uint64(chain_id), acc, addr, args[2], acc_fee, fee, ctime)
		if e != nil {
			return retErr(fmt.Errorf("create tx error: %s", e.Error()))
		}
		txbody, e := tx.Serialize()
		if e != nil {
			return retErr(fmt.Errorf("Create tx error: %s", e.Error()))
		}
		// success
		return retData(map[string]interface{}{
			"Diamonds":          diamonds.SerializeHACDlistToCommaSplitString(),
			"DiamondCount":      int(diamonds.Count),
			"Fee":               fee.ToFinString(),
			"TxHash":            tx.Hash().ToHex(),
			"TxBody":            hex.EncodeToString(txbody),
			"FeeAddress":        acc_fee.AddressReadable,
			"PaymentAddress":    acc.AddressReadable,
			"CollectionAddress": args[2],
			"Timestamp":         strconv.FormatInt(ctime, 10),
		})
	})

}
