package wasmsdk

import (
	"syscall/js"
)

func HacTransferSDK() {

	/* CreateHacTransfer */
	jsGlobalRegFunc("CreateHacTransfer", func(that js.Value, args []js.Value) interface{} {

		// success
		return retData(map[string]interface{}{
			"TxHash": "",
			"TxBody": "",
		})
	})

}

