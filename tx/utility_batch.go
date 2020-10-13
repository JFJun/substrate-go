package tx

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/JFJun/substrate-go/scale"
	"github.com/JFJun/substrate-go/types"
)

// 0x4902 84
// 4abed396aa071ba6cb7ecd96f48f9d8af2e8d5d9d5a248fd776a14307ba0782c
// 01
// 20d7a913354e6c21865754d62b08dbd933e150fd816f5bcedc6dcf3922a99631
// 60bb0a950e3f427d63e73ea72999c2fc71bb256b24bfa2edcd5f52fc4b46b48e
// d50000001a000405001cf326c5aaa5af9f0e2791e66310fe8f044faadaf12567eaa0976959d1f7731f0b009503b48415
//----------------------------------------------------------------------------------------

// 0x4902 84
// 760481bc6b19a1265ea4f7f88ac07e8130954b0a8e6007c8c8532898e6b94967
// 01
// 6cee2c7a48d854b8e4805f64ec9ff4e1d2200b668394e14fc2ff2e514b02043b
// e03e0ca5f9628ea1cc8e8a551c5b9c36f0dba31222005da7bd2105a8f4a70b89
// d50000001a000405001cf326c5aaa5af9f0e2791e66310fe8f044faadaf12567eaa0976959d1f7731f0be0a8c1677203
//----------------------------------------------------------------------------------------

// 0x4902 84
// 96ba6b66288f08964c0ff558a27e5242936ab9f8396ba288f93df932a7db9b53
// 01
// a217fdf0f98e138cde0b4488f526084c854c5dfa59f0153970e0bdd289f42156
// 6801ff64b40052913b763f0eb968919ffa3206e6ef5851a3410d5d2a9c387b8c
// d5000000
//	1a00  callidx
// 04	  数量
// 0500   callidx
//	1cf326c5aaa5af9f0e2791e66310fe8f044faadaf12567eaa0976959d1f7731f0b00ad57158c04

// 0xf102 84
// 1cf326c5aaa5af9f0e2791e66310fe8f044faadaf12567eaa0976959d1f7731f
// 01
// ccbd9cb4e8d8b79cc59b9064034b7eddc2ea54fe0c442ec070b6cbd00384621a
// e62cd7bbbe5d184098f97c59d29b5823a9c5c19d6dfa95e6cfb200053821ba89
// d5031269010000
// 1a00
// 08
// 0500
// e7c99e923e443bd57df9d4055a262aa898fd19639e22e902923bb7ba1684818e07006adb2a7b
// 0500
// 4ef03e2099425007c3204da4c528409fbc1b32c9b91f37c0c408fd6770c9355d07803342864c

type UtilityBatch struct {
	Calls []MethodTransfer
}

func NewUtilityBatch(pubkey_amount map[string]uint64) (*UtilityBatch, error) {
	if len(pubkey_amount) <= 0 {
		return nil, errors.New("params is null")
	}
	var calls []MethodTransfer
	for pubkey, amount := range pubkey_amount {
		pubBytes, err := hex.DecodeString(pubkey)
		if err != nil || len(pubBytes) != 32 {
			return nil, errors.New("invalid dest public key")
		}

		if amount == 0 {
			return nil, errors.New("zero amount")
		}
		uAmount := types.UCompact(amount)
		var buffer = bytes.Buffer{}
		s := scale.NewEncoder(&buffer)
		errA := uAmount.Encode(*s)
		if errA != nil {
			return nil, fmt.Errorf("encode amount error,Err=[%v]", errA)
		}
		calls = append(calls, MethodTransfer{
			DestPubkey: pubBytes,
			Amount:     buffer.Bytes(),
		})
	}
	return &UtilityBatch{
		Calls: calls,
	}, nil
}

func (mt UtilityBatch) ToBytes(utilityCallId, balanceTransferCallId string) ([]byte, error) {
	if len(mt.Calls) <= 0 {
		return nil, errors.New("calls is null")
	}
	ret, _ := hex.DecodeString(utilityCallId)
	//编码 Vec<Call>
	num := len(mt.Calls)
	uNum := types.UCompact(num)
	var numBuf = bytes.Buffer{}
	s := scale.NewEncoder(&numBuf)
	err := uNum.Encode(*s)
	if err != nil {
		return nil, fmt.Errorf("encode calls num error,Err=[%v]", err)
	}
	ret = append(ret, numBuf.Bytes()...)
	for _, call := range mt.Calls {
		callValue, err := call.ToBytes(balanceTransferCallId)
		if err != nil {
			return nil, fmt.Errorf("encode call value error,Err=%v", err)
		}
		ret = append(ret, callValue...)
	}
	return ret, nil
}
