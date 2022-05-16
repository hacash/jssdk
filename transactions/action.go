package transactions

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/interfaces"
)

/* *********************************************************** */

func NewActionByKind(kind uint16) (interfaces.Action, error) {
	////////////////////   ACTIONS   ////////////////////
	switch kind {
	case 1:
		return new(Action_1_SimpleToTransfer), nil
	case 4:
		return new(Action_4_DiamondCreate), nil
	case 5:
		return new(Action_5_DiamondTransfer), nil
	case 6:
		return new(Action_6_OutfeeQuantityDiamondTransfer), nil
	case 8:
		return new(Action_8_SimpleSatoshiTransfer), nil
	case 11:
		return new(Action_11_FromToSatoshiTransfer), nil
	case 13:
		return new(Action_13_FromTransfer), nil
	case 14:
		return new(Action_14_FromToTransfer), nil
	case 28:
		return new(Action_28_FromSatoshiTransfer), nil

	}
	////////////////////    END      ////////////////////
	return nil, fmt.Errorf("Cannot find Action kind of %d.", +kind)
}

func ParseAction(buf []byte, seek uint32) (interfaces.Action, uint32, error) {
	if seek+2 >= uint32(len(buf)) {
		return nil, 0, fmt.Errorf("[ParseAction] seek out of buf len.")
	}
	var kind = binary.BigEndian.Uint16(buf[seek : seek+2])
	var act, e1 = NewActionByKind(kind)
	if e1 != nil {
		return nil, 0, e1
	}
	var mv, err = act.Parse(buf, seek+2)
	return act, mv, err
}
