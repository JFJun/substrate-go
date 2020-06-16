package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JFJun/substrate-go/config"
	v11 "github.com/JFJun/substrate-go/model/v11"
	"github.com/JFJun/substrate-go/scale"
	"github.com/JFJun/substrate-go/ss58"
	"github.com/JFJun/substrate-go/state"
	"github.com/JFJun/substrate-go/util"
	codes "github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/source"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
	"golang.org/x/crypto/blake2b"
	"math/big"
	"strconv"
	"strings"
)

type Client struct {
	Rpc                *util.RpcClient
	Metadata           *codes.MetadataDecoder
	CoinType           string
	SpecVersion        int
	TransactionVersion int
	genesisHash        string
}

func New(url, user, password string) (*Client, error) {
	client := new(Client)
	if strings.HasPrefix(url, "wss") {
		//todo 连接websocket
		return client, errors.New("do not support websocket")
	}
	client.Rpc = util.New(url, user, password)
	//初始化运行版本
	err := client.initRuntimeVersion()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) initMetaData() error {
	metadataBytes, err := client.Rpc.SendRequest("state_getMetadata", []interface{}{})
	if err != nil {
		return fmt.Errorf("rpc get metadata error,err=%v", err)
	}
	metadata := string(metadataBytes)
	metadata = util.RemoveHex0x(metadata)
	data, err := hex.DecodeString(metadata)
	if err != nil {
		return err
	}
	m := codes.MetadataDecoder{}
	m.Init(data)
	if err := m.Process(); err != nil {
		return fmt.Errorf("parse metadata error,err=%v", err)
	}
	client.Metadata = &m
	return nil
}

func (client *Client) initRuntimeVersion() error {
	data, err := client.Rpc.SendRequest("state_getRuntimeVersion", []interface{}{})
	if err != nil {
		return fmt.Errorf("init runtime version error,err=%v", err)
	}
	var result map[string]interface{}
	errJ := json.Unmarshal(data, &result)
	if errJ != nil {
		return fmt.Errorf("init runtime version error,err=%v", errJ)
	}
	client.CoinType = strings.ToLower(result["specName"].(string))
	client.TransactionVersion = int(result["transactionVersion"].(float64))
	specVersion := int(result["specVersion"].(float64))
	// metadata 会动态改变，所以通过specVersion去检测metadata的改变
	if client.SpecVersion != specVersion {
		client.SpecVersion = specVersion
		return client.initMetaData()
	}
	client.SpecVersion = specVersion
	return nil
}

func (client *Client) GetBlockNumber(blockHash string) (int64, error) {
	var (
		resp []byte
		err  error
	)
	if blockHash == "" {
		blockHash, err = client.GetFinalizedHead()
		if err != nil {
			return -1, err
		}
	}
	resp, err = client.Rpc.SendRequest("chain_getBlock", []interface{}{blockHash})
	if err != nil {
		return -1, err
	}
	var block v11.SignedBlock
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return -1, err
	}

	b, isOK := new(big.Int).SetString(block.Block.Header.Number[2:], 16)
	if !isOK {
		return -1, errors.New("parse hex block number error")
	}
	return b.Int64(), nil
}

func (client *Client) GetFinalizedHead() (string, error) {
	resp, err := client.Rpc.SendRequest("chain_getFinalizedHead", []interface{}{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (client *Client) GetGenesisHash() string {
	if client.genesisHash != "" {
		return client.genesisHash
	}
	resp, err := client.Rpc.SendRequest("chain_getBlockHash", []interface{}{0})
	if err != nil {
		return ""
	}
	client.genesisHash = string(resp)
	return string(resp)
}

func (client *Client) GetAccountInfo(address string) ([]byte, error) {
	errV := client.initRuntimeVersion()
	if errV != nil {
		return nil, errV
	}
	pub, err := ss58.DecodeToPub(address)
	if err != nil {
		return nil, err
	}
	key, err1 := state.CreateStorageKey(client.Metadata, "System", "Account", pub, nil)
	if err1 != nil {
		return nil, fmt.Errorf("create stroage key error,err=%v", err1)
	}

	resp, err2 := client.Rpc.SendRequest("state_getStorageAt", []interface{}{key})
	if err2 != nil {
		return nil, err2
	}
	respStr := util.RemoveHex0x(string(resp))
	data, _ := hex.DecodeString(respStr)
	raw := state.NewStorageDataRaw(data)
	var target state.AccountInfo
	scale.NewDecoder(bytes.NewReader(raw)).Decode(&target)
	if &target == nil {
		return nil, fmt.Errorf("decode stroage data error,data=[%s]", data)
	}
	return json.Marshal(target)
}

/*
根据高度获取对应的区块信息以及交易信息
*/
func (client *Client) GetBlockByNumber(height int64) (*v11.BlockResponse, error) {
	var (
		respData []byte
		err      error
	)
	respData, err = client.Rpc.SendRequest("chain_getBlockHash", []interface{}{height})
	if err != nil || len(respData) == 0 {
		return nil, fmt.Errorf("get block hash error,err=%v", err)
	}
	blockHash := string(respData)

	return client.GetBlockByHash(blockHash)
}

func (client *Client) GetBlockByHash(blockHash string) (*v11.BlockResponse, error) {
	var (
		respData []byte
		err      error
	)
	errV := client.initRuntimeVersion()
	if errV != nil {
		return nil, errV
	}
	respData, err = client.Rpc.SendRequest("chain_getBlock", []interface{}{blockHash})
	if err != nil || len(respData) == 0 {
		return nil, fmt.Errorf("get block error,err=%v", err)
	}
	var block v11.SignedBlock
	err = json.Unmarshal(respData, &block)
	if err != nil {
		return nil, fmt.Errorf("parse block error")
	}
	blockResp := new(v11.BlockResponse)
	number, _ := strconv.ParseInt(util.RemoveHex0x(block.Block.Header.Number), 16, 64)
	blockResp.Height = number
	blockResp.ParentHash = block.Block.Header.ParentHash
	blockResp.BlockHash = blockHash
	if len(block.Block.Extrinsics) > 0 {
		//extrinsicNum:=len(block.Block.Extrinsics)
		err = client.parseExtrinsicByDecode(block.Block.Extrinsics, blockResp)
		if err != nil {
			return nil, err
		}
		err = client.parseExtrinsicByStorage(blockHash, blockResp)
		if err != nil {
			return nil, err
		}
	}

	return blockResp, nil
}

type parseBlockExtrinsicParams struct {
	from, to, sig, era, txid string
	nonce                    int64
	extrinsicIdx             int
}

func (client *Client) parseExtrinsicByDecode(extrinsics []string, blockResp *v11.BlockResponse) error {

	var (
		params    []parseBlockExtrinsicParams
		timestamp int64
		//idx int
	)
	defer func() {
		if err := recover(); err != nil {
			blockResp.Timestamp = timestamp
			blockResp.Extrinsic = []*v11.ExtrinsicResponse{}
			fmt.Printf("panic unkown decode type: err=%v\n", err)
		}
	}()

	for i, extrinsic := range extrinsics {
		//idx = i
		e := codes.ExtrinsicDecoder{}
		option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata}
		e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsic)}, &option)
		e.Process()
		bb, err := json.Marshal(e.Value)
		if err != nil {
			return fmt.Errorf("parse extrinsic error,err=%v", err)
		}
		var resp v11.ExtrinsicDecodeResponse
		errM := json.Unmarshal(bb, &resp)
		if errM != nil {
			return fmt.Errorf("json unmarshal extrinsic error,err=%v", errM)
		}
		switch resp.CallModule {
		case "Timestamp":
			for _, param := range resp.Params {
				if param.Name == "now" {
					timestamp = int64(param.Value.(float64))
				}
			}
		case "Balances":
			blockData := parseBlockExtrinsicParams{}
			blockData.from, _ = ss58.EncodeByPubHex(resp.AccountId, config.PrefixMap[client.CoinType])
			blockData.era = resp.Era
			blockData.sig = resp.Signature
			blockData.nonce = resp.Nonce
			blockData.extrinsicIdx = i
			blockData.txid = createTxHash(extrinsic)
			for _, param := range resp.Params {
				if param.Name == "dest" {
					blockData.to, _ = ss58.EncodeByPubHex(param.ValueRaw, config.PrefixMap[client.CoinType])
				}
			}

			params = append(params, blockData)
		case "Claims": //crab 转账call_module
			blockData := parseBlockExtrinsicParams{}
			blockData.from, _ = ss58.EncodeByPubHex(resp.AccountId, config.PrefixMap[client.CoinType])
			blockData.era = resp.Era
			blockData.sig = resp.Signature
			blockData.nonce = resp.Nonce
			blockData.extrinsicIdx = i
			blockData.txid = createTxHash(extrinsic)
			for _, param := range resp.Params {
				if param.Name == "dest" {
					blockData.to, _ = ss58.EncodeByPubHex(param.ValueRaw, config.PrefixMap[client.CoinType])
				}
			}
			params = append(params, blockData)
		default:
			//todo  add another call_module 币种不同可能使用的call_module不一样
			continue
		}

	}
	blockResp.Timestamp = timestamp
	//解析params
	if len(params) == 0 {

		blockResp.Extrinsic = []*v11.ExtrinsicResponse{}
		return nil
	}
	blockResp.Extrinsic = make([]*v11.ExtrinsicResponse, len(params))
	for idx, param := range params {

		e := new(v11.ExtrinsicResponse)
		e.Signature = param.sig
		e.FromAddress = param.from
		e.ToAddress = param.to
		e.Nonce = param.nonce
		e.Era = param.era
		e.ExtrinsicIndex = param.extrinsicIdx
		e.Txid = param.txid
		blockResp.Extrinsic[idx] = e
	}

	return nil
}

func (client *Client) parseExtrinsicByStorage(blockHash string, blockResp *v11.BlockResponse) error {
	var (
		err  error
		key  string
		resp []byte
	)
	key, err = state.CreateStorageKey(client.Metadata, "System", "Events", nil, nil)
	if err != nil {
		return fmt.Errorf("create stroage key error,err=%v", err)
	}
	resp, err = client.Rpc.SendRequest("state_getStorageAt", []interface{}{key, blockHash})
	if err != nil || len(resp) <= 0 {
		return fmt.Errorf("get system events error,err=%v", err)
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
		return err
	}
	var eventResp []v11.EventResponse
	err = json.Unmarshal(data, &eventResp)
	if err != nil {
		return fmt.Errorf("parse events error,err=%v", err)
	}
	if len(eventResp) > 0 {

		for _, event := range eventResp {
			var (
				defaultSuccess = "success"
				amount         = "0"
			)
			switch event.EventId {
			case config.ExtrinsicFailed:
				defaultSuccess = "failed"
				break
			case config.Transfer:
				if event.ModuleId == "Balances" {
					if len(event.Params) <= 0 {
						defaultSuccess = "failed"
						break
					}
					for _, param := range event.Params {
						if param.Type == "Balance" {
							amount = param.Value.(string)
						}
					}
				}
			default:
				continue
			}
			for _, e := range blockResp.Extrinsic {
				if e.ExtrinsicIndex == event.ExtrinsicIdx {
					e.Type = "transfer"
					e.Amount = amount
					e.Status = defaultSuccess
					e.Fee = client.calcFee(eventResp, event.ExtrinsicIdx)
				}
			}

		}
		////设置交易状态
		//blockResp.Status = defaultSuccess
		//if defaultSuccess=="failed" {
		//	return nil
		//}
		////在做一次for循环计算手续费
		//blockResp.Extrinsic.Fee=client.calcFee(eventResp,extrinsicIdx)
	}
	return err
}

/*
todo maybe have anther fee events
*/
func (client *Client) calcFee(events []v11.EventResponse, extrinsicIdx int) string {
	fee := new(big.Int).SetInt64(0)
	for _, event := range events {
		if event.ExtrinsicIdx == extrinsicIdx {
			if config.IsContainFeeEventId(event.EventId) {
				switch event.ModuleId {
				case "Treasury":
					if len(event.Params) == 0 {
						continue
					}
					for _, param := range event.Params {
						if strings.Contains(param.Type, "Balance") {
							value := param.Value.(string)
							subFee, isOk := new(big.Int).SetString(value, 10)
							if !isOk {
								continue
							}
							fee = fee.Add(fee, subFee)
						}
					}
				case "Balances":
					if len(event.Params) == 0 {
						continue
					}
					for _, param := range event.Params {
						if strings.Contains(param.Type, "Balance") {
							value := param.Value.(string)
							subFee, isOk := new(big.Int).SetString(value, 10)
							if !isOk {
								continue
							}
							fee = fee.Add(fee, subFee)
						}
					}

				default:
					continue
				}
			}
		}
	}
	return fee.String()
}

func createTxHash(extrinsic string) string {
	data, _ := hex.DecodeString(util.RemoveHex0x(extrinsic))
	d := blake2b.Sum256(data)
	return "0x" + hex.EncodeToString(d[:])
}
