package test

import (
	"encoding/json"
	"fmt"
	"github.com/JFJun/substrate-go/rpc"
	"testing"
)

func Test_GetBlock(t *testing.T) {
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	data, err := client.GetBlockByNumber(1801866)
	if err != nil {
		panic(err)
	}
	d, _ := json.Marshal(data)
	fmt.Println(string(d))
}
