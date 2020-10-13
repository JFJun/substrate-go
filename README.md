# substrate go sdk
## 注意
    请使用最新的包版本v1.2.0
    目前这个包写的比较乱，后续会对这个包进行一个大的改动，会升级到v2.0
    原因之一是因为使用的第三方包(github.com/itering/scale.go)里面存在panic，导致有些区块解析不出来，所以自己重写了一个包。后续会进行替换
## 介绍
    目前该包支持dot,ksm以及crab(其他的substrate下的币种没测，应该也能用)的以下功能,如果要解析pcx，可以使用https://github.com/JFJun/chainX-go.git
    目前主要提供了三个功能
    1. 获取区块并解析区块extrinsic
    2. 根据地址获取可用余额以及nonce
    3. 交易离线签名（目前仅支持er25519）

## 解析区块
        解析区块需要判断以下两个字段：
        1. 判断status是否为"success"，如果是"success"，表示里面的交易是有效的
        2. 判断type是否为"transfer"，如果是"transfer"，表示里面有extrinsic交易
        例如：
        var client,err =rpc2.New("url","","")
        err != nil{
            return
        }
    	resp,err:=client.Rpc.SendRequest("chain_getFinalizedHead",[]interface{}{})
    	if err != nil {
    		return
    	}
    	blockHash:=string(resp)
    	block,err:=client.GetBlockByHash(blockHash)
    	if err != nil {
    		return
    	}
    	if block.Status!="success" {		//success-> valid transaction ELSE-> invalid transaction
    		return
    	}
    	if block.Extrinsic.Type!="transfer" {	//transfer -> Balance.Transfer  ELSE-> unknown type
    		return
    	}
    	//...
## 获取账户余额
    var client,err =rpc2.New("url","","")
            err != nil{
                return
            }
    data,err:=client.GetAccountInfo("address")
    	if err != nil {
    		return
    	}
    	fmt.Println(string(data))

## 离线签名以及转账
#### 1. Balance.transfer
    注意： 转账blockHash设为genesisHash，blockNumber可以设为任意
    c,err:=rpc.New("wss://rpc.polkadot.io","","")
    	if err != nil {
    		return
    	}
    	btTx:=tx.CreateTransaction("from","to",10000000,12,0)
    	btTx.SetGenesisHashAndBlockHash("genesisHash","genesisHash",0)
    	// 通过方法去获取callIdx，不走config
    	callIdx,err:=c.GetCallIdx("Balances","transfer")
    	if err != nil {
    		return
    	}
    	btTx.SetSpecVersionAndCallId(uint32(c.SpecVersion),uint32(c.TransactionVersion),callIdx)
    	_,message,err:=btTx.CreateEmptyTransactionAndMessage()
    	if err != nil {
    		return
    	}
    	sig,err:=btTx.SignTransaction("private key",message)
    	if err != nil {
    		return
    	}
    	txHex,err:=btTx.GetSignTransaction(sig)
    	if err != nil {
    		return
    	}
    	//broadcast tx
    	txidBytes,err:=c.Rpc.SendRequest("author_submitExtrinsic",[]interface{}{txHex})
    	if err != nil {
    		return
    	}
    	txid:=string(txidBytes)
    	fmt.Println(txid)
#### 2. Utility.batch（批量转账）
    c,err:=rpc.New("wss://rpc.polkadot.io","","")
    	if err != nil {
    		return
    	}
    	address_amount:=make(map[string]uint64)
    	address_amount["to1"] = 123
    	address_amount["to2"] = 456
    	// .
    	// .
    	// .
    	ubCallIdx,err:=c.GetCallIdx("Utility","batch")
    	ubTx:=tx.CreateUtilityBatchTransaction("from",ubCallIdx,12,address_amount)
    	ubTx.SetGenesisHashAndBlockHash("genesisHash","genesisHash",0)
    	// 通过方法去获取callIdx，不走config
    	callIdx,err:=c.GetCallIdx("Balances","transfer")
    	if err != nil {
    		return
    	}
    	ubTx.SetSpecVersionAndCallId(uint32(c.SpecVersion),uint32(c.TransactionVersion),callIdx)
    	_,message,err:=ubTx.CreateEmptyTransactionAndMessage()
    	if err != nil {
    		return
    	}
    	sig,err:=ubTx.SignTransaction("private key",message)
    	if err != nil {
    		return
    	}
    	txHex,err:=ubTx.GetSignTransaction(sig)
    	if err != nil {
    		return
    	}
    	//broadcast tx
    	txidBytes,err:=c.Rpc.SendRequest("author_submitExtrinsic",[]interface{}{txHex})
    	if err != nil {
    		return
    	}
    	txid:=string(txidBytes)
    	fmt.Println(txid)