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

func SatTransferSDK() {

	/* CreateHacTransfer */
	jsGlobalRegFuncPmsString("CreateSatTransfer", func(args []string) interface{} {
		if len(args) != 6 {
			return retErr(fmt.Errorf("param num must be 6."))
		}
		chain_id, e := strconv.ParseInt(args[0], 10, 64)
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid chain id.", args[0]))
		}
		payacc := account.GetAccountByPrivateKeyOrPassword(args[1])
		addr, e := account.CheckReadableAddress(args[2])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid Hacash address.", args[2]))
		}
		satval, e := strconv.ParseInt(args[3], 10, 64)
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid transfer amount.", args[3]))
		}
		satamt := fields.Satoshi(uint64(satval))
		feeacc := account.GetAccountByPrivateKeyOrPassword(args[4])
		fee, e := fields.NewAmountFromString(args[5])
		if e != nil {
			return retErr(fmt.Errorf("'%s' not a valid fee amount.", args[5]))
		}
		// create tx
		ctime := time.Now().Unix()
		tx, e := transactions.CreateOneTxOfBTCTransfer(uint64(chain_id), payacc, addr, uint64(satamt), feeacc, fee, ctime)
		if e != nil {
			return retErr(fmt.Errorf("CreateOneTxOfBTCTransfer Error: %s", e))
		}
		txbody, e := tx.Serialize()
		if e != nil {
			return retErr(fmt.Errorf("Create tx error: %s", e.Error()))
		}
		// success
		return retData(map[string]interface{}{
			"Amount":            strconv.FormatInt(int64(satamt), 10),
			"Fee":               fee.ToFinString(),
			"TxHash":            tx.Hash().ToHex(),
			"TxBody":            hex.EncodeToString(txbody),
			"PaymentAddress":    payacc.AddressReadable,
			"CollectionAddress": args[2],
			"Timestamp":         strconv.FormatInt(ctime, 10),
		})
	})

}
