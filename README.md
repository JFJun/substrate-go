# substrate go sdk
## 介绍
    目前该包支持ksm以及crab的以下功能
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
    注意： 转账blockHash设为genesisHash，blockNumber可以设为任意
    originTx:= tx.CreateTransaction("from","to",uint64(amount),uint64(nonce),uint64(tip))
    	originTx.SetGenesisHashAndBlockHash("genesisHash","genesisHash",blockNumber)
    	originTx.SetSpecVersionAndCallId(uint32(client.SpecVersion),uint32(client.TransactionVersion),config.CallIdKusama)
    	_,message,err:=originTx.CreateEmptyTransactionAndMessage()
    	if err != nil {
    		return
    	}
    	sig,err:=originTx.SignTransaction("private Key(hex)",message)
    	if err != nil {
    		return
    	}
    	txHex,err:=originTx.GetSignTransaction(sig)
    	if err != nil {
    		return
    	}
    	txidBytes,err:=client.Rpc.SendRequest("author_submitExtrinsic",[]interface{}{txHex})
    	if err != nil {
    		panic(err)
    	}
    	txid:=string(txidBytes)
    	fmt.Println(txid)
