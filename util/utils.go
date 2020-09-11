package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/JFJun/substrate-go/scale"
	"github.com/JFJun/substrate-go/types"
	"github.com/JFJun/substrate-go/xxhash"
	"golang.org/x/crypto/blake2b"
	"hash"
	"strings"
)

func AppendBytes(data1, data2 []byte) []byte {
	if data2 == nil {
		return data1
	}
	return append(data1, data2...)
}

func SelectHash(method string) (hash.Hash, error) {
	switch method {
	case "Twox128":
		return xxhash.New128(nil), nil
	case "Blake2_256":
		return blake2b.New256(nil)
	case "Blake2_128":
		return blake2b.New(16, nil)
	case "Blake2_128Concat":
		return blake2b.New(16, nil)
	case "Twox64Concat":
		return xxhash.New64(nil), nil
	case "Identity":
		return nil, nil
	default:
		return nil, errors.New("unknown hash method")

	}

}

func RemoveHex0x(hexStr string) string {
	if strings.HasPrefix(hexStr, "0x") {
		return hexStr[2:]
	}
	return hexStr
}

func UCompactEncode(data uint64) ([]byte, error) {
	uAmount := types.UCompact(data)
	var buffer = bytes.Buffer{}

	s := scale.NewEncoder(&buffer)
	err := uAmount.Encode(*s)
	if err != nil {
		return nil, fmt.Errorf("encode UCompact error,Err=[%v]", err)
	}
	return buffer.Bytes(), nil
}
