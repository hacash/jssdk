package transactions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"math/big"
	"time"
)

type Transaction_2_Simple struct {
	Timestamp   fields.BlockTxTimestamp
	MainAddress fields.Address
	Fee         fields.Amount

	ActionCount fields.VarUint2
	Actions     []interfaces.Action

	SignCount fields.VarUint2
	Signs     []fields.Sign

	MultisignCount fields.VarUint2
	Multisigns     []fields.Multisign

	// cache data
	hashwithfee fields.Hash
	hashnofee   fields.Hash
}

func NewEmptyTransaction_2_Simple(master fields.Address) (*Transaction_2_Simple, error) {
	if !master.IsValid() {
		return nil, fmt.Errorf("Master Address is InValid ")
	}
	timeUnix := time.Now().Unix()
	return &Transaction_2_Simple{
		Timestamp:      fields.BlockTxTimestamp(uint64(timeUnix)),
		MainAddress:    master,
		Fee:            *fields.NewEmptyAmount(),
		ActionCount:    fields.VarUint2(0),
		SignCount:      fields.VarUint2(0),
		MultisignCount: fields.VarUint2(0),
	}, nil
}

func (trs *Transaction_2_Simple) Type() uint8 {
	return 2
}

func (trs *Transaction_2_Simple) ClearHash() {
	trs.hashwithfee = nil
	trs.hashnofee = nil
}

func (trs *Transaction_2_Simple) Clone() interfaces.Transaction {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_2_Simple)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
}

func (trs *Transaction_2_Simple) Copy() interfacev2.Transaction {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_2_Simple)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
}

func (trs *Transaction_2_Simple) Serialize() ([]byte, error) {
	body, e0 := trs.SerializeNoSign()
	if e0 != nil {
		return nil, e0
	}
	var buffer = new(bytes.Buffer)
	buffer.Write(body)
	// sign
	b1, e1 := trs.SignCount.Serialize()
	if e1 != nil {
		return nil, e1
	}
	buffer.Write(b1)
	for i := 0; i < int(trs.SignCount); i++ {
		var bi, e = trs.Signs[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	// muilt sign
	b2, e2 := trs.MultisignCount.Serialize()
	if e2 != nil {
		return nil, e2
	}
	buffer.Write(b2)
	for i := 0; i < int(trs.MultisignCount); i++ {
		var bi, e = trs.Multisigns[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	// ok
	return buffer.Bytes(), nil
}

func (trs *Transaction_2_Simple) SerializeNoSign() ([]byte, error) {
	return trs.SerializeNoSignEx(true)
}

// 序列化不包含签名内容的所有其它数据
func (trs *Transaction_2_Simple) SerializeNoSignEx(hasfee bool) ([]byte, error) {
	var buffer = new(bytes.Buffer)
	buffer.Write([]byte{trs.Type()}) // type
	b1, _ := trs.Timestamp.Serialize()
	buffer.Write(b1)
	b2, _ := trs.MainAddress.Serialize()
	buffer.Write(b2)
	if hasfee { // 是否需要 fee 字段
		b3, _ := trs.Fee.Serialize()
		buffer.Write(b3) // 费用付出者签名 需要fee字段， 否则不需要
	}
	b4, _ := trs.ActionCount.Serialize()
	buffer.Write(b4)
	for i := 0; i < len(trs.Actions); i++ {
		var bi, e = trs.Actions[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	//if nofee {
	//	fmt.Println( "SerializeNoSignEx: " + hex.EncodeToString(buffer.Bytes()))
	//}
	return buffer.Bytes(), nil
}

func (trs *Transaction_2_Simple) Parse(buf []byte, seek uint32) (uint32, error) {
	m1, e := trs.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	m2, e := trs.MainAddress.Parse(buf, m1)
	if e != nil {
		return 0, e
	}
	m3, e := trs.Fee.Parse(buf, m2)
	if e != nil {
		return 0, e
	}
	m4, e := trs.ActionCount.Parse(buf, m3)
	if e != nil {
		return 0, e
	}
	iseek := m4
	for i := 0; i < int(trs.ActionCount); i++ {
		var act, sk, err = ParseAction(buf, iseek)
		trs.Actions = append(trs.Actions, act)
		iseek = sk
		if err != nil {
			return 0, err
		}
	}
	iseek, e = trs.SignCount.Parse(buf, iseek)
	if e != nil {
		return 0, e
	}
	for i := 0; i < int(trs.SignCount); i++ {
		var sign fields.Sign
		iseek, e = sign.Parse(buf, iseek)
		if e != nil {
			return 0, e
		}
		trs.Signs = append(trs.Signs, sign)
	}
	iseek, e = trs.MultisignCount.Parse(buf, iseek)
	if e != nil {
		return 0, e
	}
	for i := 0; i < int(trs.MultisignCount); i++ {
		var multisign fields.Multisign
		iseek, e = multisign.Parse(buf, iseek)
		if e != nil {
			return 0, e
		}
		trs.Multisigns = append(trs.Multisigns, multisign)
	}
	return iseek, nil
}

func (trs *Transaction_2_Simple) Size() uint32 {
	totalsize := 1 +
		trs.Timestamp.Size() +
		trs.MainAddress.Size() +
		trs.Fee.Size() +
		trs.ActionCount.Size()
	for i := 0; i < int(trs.ActionCount); i++ {
		totalsize += trs.Actions[i].Size()
	}
	totalsize += trs.SignCount.Size()
	for i := 0; i < int(trs.SignCount); i++ {
		totalsize += trs.Signs[i].Size()
	}
	totalsize += trs.MultisignCount.Size()
	for i := 0; i < int(trs.MultisignCount); i++ {
		totalsize += trs.Multisigns[i].Size()
	}
	return totalsize
}

// 交易唯一哈希值
func (trs *Transaction_2_Simple) HashWithFee() fields.Hash {
	if trs.hashwithfee == nil {
		return trs.HashWithFeeFresh()
	}
	return trs.hashwithfee
}

func (trs *Transaction_2_Simple) HashWithFeeFresh() fields.Hash {
	stuff, _ := trs.SerializeNoSign()
	digest := fields.CalculateHash(stuff)
	trs.hashwithfee = digest // 缓存
	return trs.hashwithfee
}

func (trs *Transaction_2_Simple) Hash() fields.Hash {
	if trs.hashnofee == nil {
		return trs.HashFresh()
	}
	return trs.hashnofee
}

func (trs *Transaction_2_Simple) HashFresh() fields.Hash {
	is_has_fee := false
	stuff, _ := trs.SerializeNoSignEx(is_has_fee)
	digest := fields.CalculateHash(stuff)
	trs.hashnofee = digest
	return trs.hashnofee
}

func (trs *Transaction_2_Simple) AppendAction(action interfacev2.Action) error {
	if trs.ActionCount >= 65530 {
		return fmt.Errorf("Actions too much")
	}
	if trs.Actions == nil {
		trs.ActionCount = 0 // 初始化
		trs.Actions = make([]interfaces.Action, 0)
	}
	trs.ActionCount += 1
	trs.Actions = append(trs.Actions, action.(interfaces.Action))
	trs.ClearHash() // 重置哈希缓存
	return nil
}

func (trs *Transaction_2_Simple) AddAction(action interfaces.Action) error {
	if trs.ActionCount >= 65530 {
		return fmt.Errorf("Actions too much")
	}
	if trs.Actions == nil {
		trs.ActionCount = 0 // 初始化
		trs.Actions = make([]interfaces.Action, 0)
	}
	trs.ActionCount += 1
	trs.Actions = append(trs.Actions, action)
	trs.ClearHash() // 重置哈希缓存
	return nil
}

// 从 actions 拿出需要签名的地址
func (trs *Transaction_2_Simple) RequestSignAddresses(reqs []fields.Address, dropfeeaddr bool) ([]fields.Address, error) {
	if !trs.MainAddress.IsValid() {
		return nil, fmt.Errorf("Master Address is InValid ")
	}
	requests := make([]fields.Address, 0, 32)
	// 另外新加的需要验证的
	if reqs != nil {
		for _, r := range reqs {
			requests = append(requests, r)
		}
	}
	// 拿出 actions 的需要签名
	for i := 0; i < int(trs.ActionCount); i++ {
		actreqs := trs.Actions[i].RequestSignAddresses()
		requests = append(requests, actreqs...)
	}
	// 去重
	results := make([]fields.Address, 0, len(requests))
	has := make(map[string]bool)
	if !dropfeeaddr {
		// 不去掉，加上主地址
		results = append(results, trs.MainAddress)
	}
	// 费用方/主地址  去重
	has[string(trs.MainAddress)] = true
	for i := 0; i < len(requests); i++ {
		strkey := string(requests[i])
		if _, ok := has[strkey]; !ok {
			results = append(results, requests[i])
			has[strkey] = true // 标记重复
		}
	}
	// 返回
	return results, nil
}

// 清清除所有签名
func (trs *Transaction_2_Simple) CleanSigns() {
	trs.SignCount = 0
	trs.Signs = []fields.Sign{}
}

// 返回所有签名
func (trs *Transaction_2_Simple) GetSigns() []fields.Sign {
	return trs.Signs
}

// 设置签名数据
func (trs *Transaction_2_Simple) SetSigns(allsigns []fields.Sign) {
	num := len(allsigns)
	if num > 65535 {
		panic("Sign is too much.")
	}
	trs.SignCount = fields.VarUint2(num)
	trs.Signs = make([]fields.Sign, 0)
	trs.Signs = append(trs.Signs, allsigns...) // copy
}

// 填充单个需要的签名
func (trs *Transaction_2_Simple) FillTargetSign(signacc *account.Account) error {
	signaddr := fields.Address(signacc.Address)
	tarhash := trs.Hash()
	if signaddr.Equal(trs.MainAddress) {
		tarhash = trs.HashWithFee() // 主地址使用hash不同
	}
	addrPrivateKeys := map[string][]byte{}
	addrPrivateKeys[string(signacc.Address)] = signacc.PrivateKey
	// 执行一个签名
	return trs.addOneSign(tarhash, addrPrivateKeys, signacc.Address)
}

// 填充全部需要的签名
func (trs *Transaction_2_Simple) FillNeedSigns(addrPrivateKeys map[string][]byte, appendReqs []fields.Address) error {
	hashWithFee := trs.HashWithFee()
	hashNoFee := trs.Hash()
	requests, e0 := trs.RequestSignAddresses(appendReqs, true)
	if e0 != nil {
		return e0
	}
	// 主签名（包括手续费）
	e1 := trs.addOneSign(hashWithFee, addrPrivateKeys, trs.MainAddress)
	if e1 != nil {
		return e1
	}
	// 其他签名（不包括手续费字段）
	for i := 0; i < len(requests); i++ {
		e1 := trs.addOneSign(hashNoFee, addrPrivateKeys, requests[i])
		if e1 != nil {
			return e1
		}
	}
	// 填充成功
	return nil
}

func (trs *Transaction_2_Simple) addOneSign(hash []byte, addrPrivates map[string][]byte, address []byte) error {
	// 判断私钥是否存在
	privitebytes, has := addrPrivates[string(address)]
	if !has {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' necessary")
	}
	privite, e1 := account.GetAccountByPriviteKey(privitebytes)
	if e1 != nil {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' error")
	}
	// 判断签名是否已经存在，如果存在则去掉重新加入
	var alreadly = -1
	for i, sig := range trs.Signs {
		if bytes.Compare(sig.PublicKey, privite.PublicKey) == 0 {
			alreadly = i
			break
		}
	}
	// 计算签名
	signature, e2 := privite.Private.Sign(hash)
	if e2 != nil {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' do sign error")
	}
	sigObjSave := fields.Sign{
		PublicKey: privite.PublicKey,
		Signature: signature.Serialize64(),
	}
	if alreadly > -1 {
		// replace
		trs.Signs[alreadly] = sigObjSave
	} else {
		// append
		trs.SignCount += 1
		trs.Signs = append(trs.Signs, sigObjSave)
	}

	//// test ////
	//verok := signature.Verify(hash, privite.Private.PubKey())
	//if !verok {
	//	panic("false")
	//}

	return nil
}

// 单独验证其中一个签名
func (trs *Transaction_2_Simple) VerifyTargetSigns(reqaddrs []fields.Address) (bool, error) {
	otherhash := trs.Hash()
	mainhash := trs.HashWithFee()
	// 全部签名
	allSigns := make(map[string]fields.Sign)
	for i := 0; i < len(trs.Signs); i++ {
		sig := trs.Signs[i]
		addrbts := account.NewAddressFromPublicKeyV0(sig.PublicKey)
		addr := fields.Address(addrbts)
		allSigns[string(addr)] = sig
	}
	// 依次验证
	for _, v := range reqaddrs {
		// 判断是否为主地址
		tarhash := otherhash // 交易哈希
		isMainAddr := v.Equal(trs.MainAddress)
		if isMainAddr { // 是否为主地址
			tarhash = mainhash
		}
		ok, e := verifyOneSignature(allSigns, v, tarhash)
		if !ok || e != nil {
			return ok, e // 验证失败
		}
		// next
	}
	// 验证成功
	return true, nil
}

// 验证需要的签名
// reqs 附加的另外要验证的
func (trs *Transaction_2_Simple) VerifyAllNeedSigns() (bool, error) {
	hashWithFee := trs.HashWithFee()
	hashNoFee := trs.Hash()
	// 开始判断
	allSigns := make(map[string]fields.Sign)
	for i := 0; i < len(trs.Signs); i++ {
		sig := trs.Signs[i]
		addr := account.NewAddressFromPublicKeyV0(sig.PublicKey)
		allSigns[string(addr)] = sig
	}
	// 验证主签名（包括手续费）
	ok, e := verifyOneSignature(allSigns, trs.MainAddress, hashWithFee)
	if e != nil || !ok {
		return ok, e
	}
	// 验证全部需要验证的签名 // 去掉主地址
	requests, e := trs.RequestSignAddresses(nil, true)
	if e != nil {
		return false, e
	}
	if requests == nil || len(requests) == 0 {
		return true, nil // 没有其他需要验证
	}
	// 验证其他所有签名（不包含手续费字段）
	for i := 0; i < len(requests); i++ {
		ok, e := verifyOneSignature(allSigns, requests[i], hashNoFee)
		if e != nil || !ok {
			return ok, e
		}
	}
	// 验证成功
	return true, nil
}

func verifyOneSignature(allSigns map[string]fields.Sign, address fields.Address, hash []byte) (bool, error) {

	main, ok := allSigns[string(address)]
	if !ok {
		return false, fmt.Errorf("address %s signature not find!", address.ToReadable())
	}
	// 检查签名
	return account.CheckSignByHash32(hash, main.PublicKey, main.Signature)
}

// 需要的余额检查
func (trs *Transaction_2_Simple) RequestAddressBalance() ([][]byte, []big.Int, error) {
	return nil, nil, nil
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_2_Simple) WriteInChainState(state interfaces.ChainStateOperation) error {
	panic("never call in transactions!")
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_2_Simple) WriteinChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

func (trs *Transaction_2_Simple) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("never call in transactions!")
}

// 手续费含量 每byte的含有多少烁代币
func (trs *Transaction_2_Simple) FeePurity() uint64 {
	panic("never call in transactions!")
}

// 查询
func (trs *Transaction_2_Simple) GetAddress() fields.Address {
	return trs.MainAddress
}

func (trs *Transaction_2_Simple) SetAddress(addr fields.Address) {
	trs.MainAddress = addr
	trs.ClearHash() // 重置哈希缓存
}

func (trs *Transaction_2_Simple) GetFeeOfMinerRealReceived() *fields.Amount {
	for _, act := range trs.Actions {
		if act.IsBurning90PersentTxFees() {
			// 销毁 90% 的tx费用
			minerReceivedFee := trs.Fee.Copy()
			if minerReceivedFee.Unit > 0 {
				// 单位下降一位（例如248变247），大小变为原来的 10%， 而销毁了 90% 。
				minerReceivedFee.Unit -= 1
			}
			// 返回矿工真实收到的竞价费，为原来的 90%
			return minerReceivedFee
		}
	}

	return &trs.Fee
}

func (trs *Transaction_2_Simple) GetFee() *fields.Amount {
	return &trs.Fee
}

func (trs *Transaction_2_Simple) SetFee(fee *fields.Amount) {
	trs.Fee = *fee
	trs.ClearHash() // 重置哈希缓存
}

func (trs *Transaction_2_Simple) GetActions() []interfacev2.Action {
	var list = make([]interfacev2.Action, len(trs.Actions))
	for i, v := range trs.Actions {
		list[i] = v.(interfacev2.Action)
	}
	return list
}

func (trs *Transaction_2_Simple) GetActionList() []interfaces.Action {
	return trs.Actions
}

func (trs *Transaction_2_Simple) GetTimestamp() uint64 { // 时间戳
	return uint64(trs.Timestamp)
}

func (trs *Transaction_2_Simple) SetMessage(fields.TrimString16) {
}

func (trs *Transaction_2_Simple) GetMessage() fields.TrimString16 {
	return fields.TrimString16("")
}
