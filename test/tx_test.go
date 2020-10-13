package test

import (
	"fmt"
	"github.com/JFJun/substrate-go/rpc"
	"github.com/JFJun/substrate-go/tx"
	"testing"
)

func Test_BalanceTransfer(t *testing.T) {
	c, err := rpc.New("wss://rpc.polkadot.io", "", "")
	if err != nil {
		return
	}
	btTx := tx.CreateTransaction("from", "to", 10000000, 12, 0)
	btTx.SetGenesisHashAndBlockHash("genesisHash", "genesisHash", 0)
	// 通过方法去获取callIdx，不走config
	callIdx, err := c.GetCallIdx("Balances", "transfer")
	if err != nil {
		return
	}
	btTx.SetSpecVersionAndCallId(uint32(c.SpecVersion), uint32(c.TransactionVersion), callIdx)
	_, message, err := btTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		return
	}
	sig, err := btTx.SignTransaction("private key", message)
	if err != nil {
		return
	}
	txHex, err := btTx.GetSignTransaction(sig)
	if err != nil {
		return
	}
	//broadcast tx
	txidBytes, err := c.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txHex})
	if err != nil {
		return
	}
	txid := string(txidBytes)
	fmt.Println(txid)
}

func Test_UtilityBatch(t *testing.T) {
	c, err := rpc.New("wss://rpc.polkadot.io", "", "")
	if err != nil {
		return
	}
	address_amount := make(map[string]uint64)
	address_amount["to1"] = 123
	address_amount["to2"] = 456
	// .
	// .
	// .
	ubCallIdx, err := c.GetCallIdx("Utility", "batch")
	ubTx := tx.CreateUtilityBatchTransaction("from", ubCallIdx, 12, address_amount)
	ubTx.SetGenesisHashAndBlockHash("genesisHash", "genesisHash", 0)
	// 通过方法去获取callIdx，不走config
	callIdx, err := c.GetCallIdx("Balances", "transfer")
	if err != nil {
		return
	}
	ubTx.SetSpecVersionAndCallId(uint32(c.SpecVersion), uint32(c.TransactionVersion), callIdx)
	_, message, err := ubTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		return
	}
	sig, err := ubTx.SignTransaction("private key", message)
	if err != nil {
		return
	}
	txHex, err := ubTx.GetSignTransaction(sig)
	if err != nil {
		return
	}
	//broadcast tx
	txidBytes, err := c.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txHex})
	if err != nil {
		return
	}
	txid := string(txidBytes)
	fmt.Println(txid)
}
