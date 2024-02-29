package transactions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"strings"
)

/**
 * 钻石交易类型
 */

// 第 20001 个钻石开始，启用 32 位的 msg byte
const DiamondCreateCustomMessageAboveNumber uint32 = 20000

// 第 30001 个钻石开始，销毁 90% 的竞价费用
const DiamondCreateBurning90PercentTxFeesAboveNumber uint32 = 30000

// 采用 30001 ~ 40000 枚钻石平均竞价费用，之前的设定为 10 枚
const DiamondStatisticsAverageBiddingBurningPriceAboveNumber uint32 = 40000

// 第 40001 个钻石，开始用 sha3_hash(diamondreshash + blockhash) 决定钻石形状和配色
const DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber uint32 = 40000

// 第 41001 个钻石，开始用 sha3_hash(diamondreshash + blockhash + bidfee) 包括竞价费参与决定钻石形状配色
const DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber uint32 = 41000

// 挖出钻石
type Action_4_DiamondCreate struct {
	Diamond  fields.DiamondName   // 钻石字面量 WTYUIAHXVMEKBSZN
	Number   fields.DiamondNumber // 钻石序号，用于难度检查
	PrevHash fields.Hash          // 上一个包含钻石的区块hash
	Nonce    fields.Bytes8        // 随机数
	Address  fields.Address       // 所属账户
	// 客户消息
	CustomMessage fields.Bytes32

	// 所属交易
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_4_DiamondCreate) Kind() uint16 {
	return 4
}

// json api
func (elm *Action_4_DiamondCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_4_DiamondCreate) Size() uint32 {
	size := 2 +
		elm.Diamond.Size() +
		elm.Number.Size() +
		elm.PrevHash.Size() +
		elm.Nonce.Size() +
		elm.Address.Size()
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		size += elm.CustomMessage.Size()
	}
	return size
}

func (elm *Action_4_DiamondCreate) GetRealCustomMessage() []byte {
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		var msgBytes, _ = elm.CustomMessage.Serialize()
		return msgBytes
	}
	return []byte{}
}

func (elm *Action_4_DiamondCreate) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var diamondBytes, _ = elm.Diamond.Serialize()
	var numberBytes, _ = elm.Number.Serialize()
	var prevBytes, _ = elm.PrevHash.Serialize()
	var nonceBytes, _ = elm.Nonce.Serialize()
	var addrBytes, _ = elm.Address.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(diamondBytes)
	buffer.Write(numberBytes)
	buffer.Write(prevBytes)
	buffer.Write(nonceBytes)
	buffer.Write(addrBytes)
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		var msgBytes, _ = elm.CustomMessage.Serialize()
		buffer.Write(msgBytes)
	}
	return buffer.Bytes(), nil
}

func (elm *Action_4_DiamondCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	moveseek1, e := elm.Diamond.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	moveseek2, e := elm.Number.Parse(buf, moveseek1)
	if e != nil {
		return 0, e
	}
	moveseek3, e := elm.PrevHash.Parse(buf, moveseek2)
	if e != nil {
		return 0, e
	}
	moveseek4, e := elm.Nonce.Parse(buf, moveseek3)
	if e != nil {
		return 0, e
	}
	moveseek5, e := elm.Address.Parse(buf, moveseek4)
	if e != nil {
		return 0, e
	}
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		moveseek5, e = elm.CustomMessage.Parse(buf, moveseek5)
		if e != nil {
			return 0, e
		}
	}
	return moveseek5, nil
}

func (elm *Action_4_DiamondCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // no sign
}

func (act *Action_4_DiamondCreate) WriteInChainState(state interfaces.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_4_DiamondCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_4_DiamondCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (elm *Action_4_DiamondCreate) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_4_DiamondCreate) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_4_DiamondCreate) IsBurning90PersentTxFees() bool {
	if uint32(act.Number) > DiamondCreateBurning90PercentTxFeesAboveNumber {
		// 从第 30001 钻石开始，销毁本笔交易的 90% 的费用
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////

// 计算钻石的可视化基因
func calculateVisualGeneByDiamondStuffHash(belong_trs interfacev2.Transaction, number uint32, stuffhx []byte, diamondstr string, peddingblkhash []byte) (fields.Bytes10, error) {
	if len(stuffhx) != 32 || len(peddingblkhash) != 32 {
		return nil, fmt.Errorf("stuffhx and peddingblkhash length must 32")
	}
	if len(diamondstr) != 16 {
		return nil, fmt.Errorf("diamondstr length must 16")
	}
	vgenehash := make([]byte, 32)
	copy(vgenehash, stuffhx)
	if number > DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber {
		// 第 40001 个钻石，开始用 sha3_hash(diamondreshash, blockhash) 决定钻石形状和配色
		vgenestuff := bytes.NewBuffer(stuffhx)
		vgenestuff.Write(peddingblkhash)
		if number > DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber {
			bidfeebts, e := belong_trs.GetFee().Serialize() // 竞价手续费
			if e != nil {
				return nil, e // 返回错误
			}
			vgenestuff.Write(bidfeebts) // 竞价费参与决定钻石形状和配色
		}
		vgenehash = fields.CalculateHash(vgenestuff.Bytes()) // 开盲盒
		// 跟区块哈希一样是随机的，需要等待钻石确认的那一刻才能知晓形状和配色
		// fmt.Println(hex.EncodeToString(vgenestuff.Bytes()))
	}
	// fmt.Printf("Calculate Visual Gene #%d, vgenehash: %s, stuffhx: %s, peddingblkhash: %s\n", number, hex.EncodeToString(vgenehash), hex.EncodeToString(stuffhx), hex.EncodeToString(peddingblkhash))

	genehexstr := make([]string, 18)
	// 前6位
	k := 0
	for i := 10; i < 16; i++ {
		s := diamondstr[i]
		e := "0"
		switch s {
		case 'W': // WTYUIAHXVMEKBSZN
			e = "0"
		case 'T':
			e = "1"
		case 'Y':
			e = "2"
		case 'U':
			e = "3"
		case 'I':
			e = "4"
		case 'A':
			e = "5"
		case 'H':
			e = "6"
		case 'X':
			e = "7"
		case 'V':
			e = "8"
		case 'M':
			e = "9"
		case 'E':
			e = "A"
		case 'K':
			e = "B"
		case 'B':
			e = "C"
		case 'S':
			e = "D"
		case 'Z':
			e = "E"
		case 'N':
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 后11位
	for i := 20; i < 31; i++ {
		x := vgenehash[i]
		x = x % 16
		e := "0"
		switch x {
		case 0:
			e = "0"
		case 1:
			e = "1"
		case 2:
			e = "2"
		case 3:
			e = "3"
		case 4:
			e = "4"
		case 5:
			e = "5"
		case 6:
			e = "6"
		case 7:
			e = "7"
		case 8:
			e = "8"
		case 9:
			e = "9"
		case 10:
			e = "A"
		case 11:
			e = "B"
		case 12:
			e = "C"
		case 13:
			e = "D"
		case 14:
			e = "E"
		case 15:
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 补齐最后一位
	genehexstr[17] = "0"
	resbts, e1 := hex.DecodeString(strings.Join(genehexstr, ""))
	if e1 != nil {
		return nil, e1
	}
	// 哈希的最后一位作为形状选择
	resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
	resbuf.Write(resbts) // 颜色选择器
	return resbuf.Bytes(), nil
}

// 计算钻石的可视化基因
func calculateVisualGeneByDiamondStuffHashV3(belong_trs interfaces.Transaction, number uint32, stuffhx []byte, diamondstr string, peddingblkhash []byte) (fields.Bytes10, error) {
	if len(stuffhx) != 32 || len(peddingblkhash) != 32 {
		return nil, fmt.Errorf("stuffhx and peddingblkhash length must 32")
	}
	if len(diamondstr) != 16 {
		return nil, fmt.Errorf("diamondstr length must 16")
	}
	vgenehash := make([]byte, 32)
	copy(vgenehash, stuffhx)
	if number > DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber {
		// 第 40001 个钻石，开始用 sha3_hash(diamondreshash, blockhash) 决定钻石形状和配色
		vgenestuff := bytes.NewBuffer(stuffhx)
		vgenestuff.Write(peddingblkhash)
		if number > DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber {
			bidfeebts, e := belong_trs.GetFee().Serialize() // 竞价手续费
			if e != nil {
				return nil, e // 返回错误
			}
			vgenestuff.Write(bidfeebts) // 竞价费参与决定钻石形状和配色
		}
		vgenehash = fields.CalculateHash(vgenestuff.Bytes()) // 开盲盒
		// 跟区块哈希一样是随机的，需要等待钻石确认的那一刻才能知晓形状和配色
		// fmt.Println(hex.EncodeToString(vgenestuff.Bytes()))
	}
	// fmt.Printf("Calculate Visual Gene #%d, vgenehash: %s, stuffhx: %s, peddingblkhash: %s\n", number, hex.EncodeToString(vgenehash), hex.EncodeToString(stuffhx), hex.EncodeToString(peddingblkhash))

	genehexstr := make([]string, 18)
	// 前6位
	k := 0
	for i := 10; i < 16; i++ {
		s := diamondstr[i]
		e := "0"
		switch s {
		case 'W': // WTYUIAHXVMEKBSZN
			e = "0"
		case 'T':
			e = "1"
		case 'Y':
			e = "2"
		case 'U':
			e = "3"
		case 'I':
			e = "4"
		case 'A':
			e = "5"
		case 'H':
			e = "6"
		case 'X':
			e = "7"
		case 'V':
			e = "8"
		case 'M':
			e = "9"
		case 'E':
			e = "A"
		case 'K':
			e = "B"
		case 'B':
			e = "C"
		case 'S':
			e = "D"
		case 'Z':
			e = "E"
		case 'N':
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 后11位
	for i := 20; i < 31; i++ {
		x := vgenehash[i]
		x = x % 16
		e := "0"
		switch x {
		case 0:
			e = "0"
		case 1:
			e = "1"
		case 2:
			e = "2"
		case 3:
			e = "3"
		case 4:
			e = "4"
		case 5:
			e = "5"
		case 6:
			e = "6"
		case 7:
			e = "7"
		case 8:
			e = "8"
		case 9:
			e = "9"
		case 10:
			e = "A"
		case 11:
			e = "B"
		case 12:
			e = "C"
		case 13:
			e = "D"
		case 14:
			e = "E"
		case 15:
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 补齐最后一位
	genehexstr[17] = "0"
	resbts, e1 := hex.DecodeString(strings.Join(genehexstr, ""))
	if e1 != nil {
		return nil, e1
	}
	// 哈希的最后一位作为形状选择
	resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
	resbuf.Write(resbts) // 颜色选择器
	return resbuf.Bytes(), nil
}

///////////////////////////////////////////////////////////////

// 转移钻石
type Action_5_DiamondTransfer struct {
	Diamond   fields.DiamondName // 钻石字面量 WTYUIAHXVMEKBSZN
	ToAddress fields.Address     // 收钻方账户

	// 数据指针
	// 所属交易
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_5_DiamondTransfer) Kind() uint16 {
	return 5
}

// json api
func (elm *Action_5_DiamondTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_5_DiamondTransfer) Size() uint32 {
	return 2 + elm.Diamond.Size() + elm.ToAddress.Size()
}

func (elm *Action_5_DiamondTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var diamondBytes, _ = elm.Diamond.Serialize()
	var addrBytes, _ = elm.ToAddress.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(diamondBytes)
	buffer.Write(addrBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_5_DiamondTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var moveseek1, _ = elm.Diamond.Parse(buf, seek)
	var moveseek2, _ = elm.ToAddress.Parse(buf, moveseek1)
	return moveseek2, nil
}

func (elm *Action_5_DiamondTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_5_DiamondTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_5_DiamondTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_5_DiamondTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (elm *Action_5_DiamondTransfer) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_5_DiamondTransfer) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_5_DiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////

// 批量转移钻石
type Action_6_OutfeeQuantityDiamondTransfer struct {
	FromAddress fields.Address              // 拥有钻石的账户
	ToAddress   fields.Address              // 收钻方账户
	DiamondList fields.DiamondListMaxLen200 // 钻石列表

	// 数据指针
	// 所属交易
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Kind() uint16 {
	return 6
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Size() uint32 {
	return 2 +
		elm.FromAddress.Size() +
		elm.ToAddress.Size() +
		elm.DiamondList.Size() // 每个钻石长6位
}

// json api
func (elm *Action_6_OutfeeQuantityDiamondTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addr1Bytes, _ = elm.FromAddress.Serialize()
	var addr2Bytes, _ = elm.ToAddress.Serialize()
	var diaBytes, e = elm.DiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addr1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(diaBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.FromAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.DiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 1) // 需from签名
	reqs[0] = elm.FromAddress
	return reqs
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_6_OutfeeQuantityDiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

// 获取区块钻石的名称列表
func (elm *Action_6_OutfeeQuantityDiamondTransfer) GetDiamondNamesSplitByComma() string {
	return elm.DiamondList.SerializeHACDlistToCommaSplitString()
}

///////////////////////////////////////////////////////////////////////

// Bulk transfer of diamonds
type Action_7_MultipleDiamondTransfer struct {
	ToAddress   fields.Address              // receive address
	DiamondList fields.DiamondListMaxLen200 // Diamond list

	// Data pointer
	// Transaction
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_7_MultipleDiamondTransfer) Kind() uint16 {
	return 7
}

func (elm *Action_7_MultipleDiamondTransfer) Size() uint32 {
	return 2 +
		elm.ToAddress.Size() +
		elm.DiamondList.Size() // Each diamond is 6 digits long
}

// json api
func (elm *Action_7_MultipleDiamondTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_7_MultipleDiamondTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addrBytes, _ = elm.ToAddress.Serialize()
	var diaBytes, e = elm.DiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addrBytes)
	buffer.Write(diaBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_7_MultipleDiamondTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.ToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.DiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_7_MultipleDiamondTransfer) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 0) // ed from address sign
	return reqs
}

func (act *Action_7_MultipleDiamondTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_7_MultipleDiamondTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (act *Action_7_MultipleDiamondTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (elm *Action_7_MultipleDiamondTransfer) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_7_MultipleDiamondTransfer) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // IsBurning 90 PersentTxFees
func (act *Action_7_MultipleDiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

// Get the name list of block diamonds
func (elm *Action_7_MultipleDiamondTransfer) GetDiamondNamesSplitByComma() string {
	return elm.DiamondList.SerializeHACDlistToCommaSplitString()
}
