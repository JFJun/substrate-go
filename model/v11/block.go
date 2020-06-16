package v11

import (
	"encoding/json"
	"github.com/JFJun/substrate-go/util"
	"strconv"
	"strings"
)

type SignedBlock struct {
	Block         Block `json:"block"`
	Justification Bytes `json:"justification"`
}

type Block struct {
	Extrinsics []string `json:"extrinsics"`
	Header     Header   `json:"header"`
}

type Header struct {
	ParentHash     string `json:"parentHash"`
	Number         string `json:"number"`
	StateRoot      string `json:"stateRoot"`
	ExtrinsicsRoot string `json:"extrinsicsRoot"`
	//Digest         interface{}   `json:"digest"`
}

type BlockNumber util.U32

// UnmarshalJSON fills BlockNumber with the JSON encoded byte array given by bz
func (b *BlockNumber) UnmarshalJSON(bz []byte) error {
	var tmp string
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}

	s := strings.TrimPrefix(tmp, "0x")

	p, err := strconv.ParseUint(s, 16, 32)
	*b = BlockNumber(p)
	return err
}

// MarshalJSON returns a JSON encoded byte array of BlockNumber
func (b BlockNumber) MarshalJSON() ([]byte, error) {
	s := strconv.FormatUint(uint64(b), 16)
	return json.Marshal(s)
}

// Encode implements encoding for BlockNumber, which just unwraps the bytes of BlockNumber
func (b BlockNumber) Encode(encoder util.Encoder) error {
	return encoder.EncodeUintCompact(uint64(b))
}

// Decode implements decoding for BlockNumber, which just wraps the bytes in BlockNumber
func (b *BlockNumber) Decode(decoder util.Decoder) error {
	u, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}
	*b = BlockNumber(u)
	return err
}

//===================================response block

type BlockResponse struct {
	Height     int64                `json:"height"`
	ParentHash string               `json:"parent_hash"`
	BlockHash  string               `json:"block_hash"`
	Timestamp  int64                `json:"timestamp"`
	Extrinsic  []*ExtrinsicResponse `json:"extrinsic"`
}

type ExtrinsicResponse struct {
	Type           string `json:"type"`   //Transfer or another
	Status         string `json:"status"` //success or fail
	FromAddress    string `json:"from_address"`
	ToAddress      string `json:"to_address"`
	Amount         string `json:"amount"`
	Fee            string `json:"fee"`
	Signature      string `json:"signature"`
	Nonce          int64  `json:"nonce"`
	Era            string `json:"era"`
	ExtrinsicIndex int    `json:"extrinsic_index"`
}
