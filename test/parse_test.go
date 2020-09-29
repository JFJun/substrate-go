package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/JFJun/substrate-go/config"
	v11 "github.com/JFJun/substrate-go/model/v11"
	"github.com/JFJun/substrate-go/rpc"
	"github.com/JFJun/substrate-go/state"
	codes "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"testing"
)

func Test_ParseExtrinsic(t *testing.T) {
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	extrinsic := "0x39028492e0feb85e225ee7c1800966ab1e69d2e0afe8021304421f9ca9bf5ea9f9784601b418c283bdadb2243dedb254d46af218c82e51048600432d473c2d0c6621d57860867599803edcc62c2bb57967dbb14b44c096bae7125142e8518be9e17e418df50200000500704f74890129ae1e780a1dcaa27fd395bfce5744d4e377dd197174537df36702070010a5d4e8"
	e := codes.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsic)}, &option)
	e.Process()
	bb, err := json.Marshal(e.Value)
	if err != nil {
		panic(err)
	}
	var resp v11.ExtrinsicDecodeResponse
	err = json.Unmarshal(bb, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))
}
func Test_parseEvent(t *testing.T) {
	var (
		err  error
		key  string
		resp []byte
	)
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	blockHash := "0xcbdb8536a723abfedad1faa70e845da65b579260347e2681b64f7eff8619a0fe"
	key, err = state.CreateStorageKey(client.Metadata, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Rpc.SendRequest("state_getStorageAt", []interface{}{key, blockHash})
	if err != nil || len(resp) <= 0 {
		panic(err)
	}
	eventsHex := string(resp)
	//解析events
	option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata, Spec: client.SpecVersion}
	ccHex := config.CoinEventType[client.CoinType]
	cc, _ := hex.DecodeString(ccHex)
	types.RegCustomTypes(source.LoadTypeRegistry(cc))
	e := codes.EventsDecoder{}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventsHex)}, &option)
	e.Process()
	data, err1 := json.Marshal(e.Value)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(data))
}
