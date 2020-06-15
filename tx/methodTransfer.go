package tx

import (
	"encoding/hex"
	"errors"
	codec "github.com/JFJun/substrate-go/codes"
)


const (
	SigningBitV4 = byte(0x84)
	Compact_U32 = "Compact<u32>"
	AccounntIDFollow = false
)

type MethodTransfer struct {
	DestPubkey []byte
	Amount     []byte
}

func NewMethodTransfer(pubkey string, amount uint64) (*MethodTransfer, error) {
	pubBytes, err := hex.DecodeString(pubkey)
	if  err != nil || len(pubBytes) != 32 {
		return nil, errors.New("invalid dest public key")
	}

	if amount == 0 {
		return nil, errors.New("zero amount")
	}
	amountStr, err := codec.Encode("Compact<u32>", amount)
	if err != nil {
		return nil, errors.New("invalid amount")
	}

	amountBytes, _ := hex.DecodeString(amountStr)

	return &MethodTransfer{
		DestPubkey: pubBytes,
		Amount:     amountBytes,
	}, nil
}

func (mt MethodTransfer) ToBytes(callId string) ([]byte, error) {

	if mt.DestPubkey == nil || len(mt.DestPubkey) != 32 || mt.Amount == nil || len(mt.Amount) == 0 {
		return nil, errors.New("invalid method")
	}

	ret, _ := hex.DecodeString(callId)
	if AccounntIDFollow {
		ret = append(ret, 0xff)
	}

	ret = append(ret, mt.DestPubkey...)
	ret = append(ret, mt.Amount...)

	return ret, nil
}
