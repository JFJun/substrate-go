package types

import (
	"github.com/JFJun/substrate-go/scale"
)

// TODO adjust to use U256 or even big ints instead, needs to adopt codec though
type UCompact uint64

func (u *UCompact) Decode(decoder scale.Decoder) error {
	ui, err := decoder.DecodeUintCompact()
	if err != nil {
		//fmt.Println("UCompact: ",434343434)
		return err
	}

	*u = UCompact(ui)
	return nil
}

func (u UCompact) Encode(encoder scale.Encoder) error {
	err := encoder.EncodeUintCompact(uint64(u))
	if err != nil {
		return err
	}
	return nil
}
