package util

import (
	"encoding/json"
	"fmt"
	//"github.com/centrifuge/go-substrate-rpc-client/scale"

	"math/big"
)

// U8 is an unsigned 8-bit integer
type U8 uint8

// NewU8 creates a new U8 type
func NewU8(u uint8) U8 {
	return U8(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U8) UnmarshalJSON(b []byte) error {
	var tmp uint8
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U8(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u U8) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(u))
}

// U16 is an unsigned 16-bit integer
type U16 uint16

// NewU16 creates a new U16 type
func NewU16(u uint16) U16 {
	return U16(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U16) UnmarshalJSON(b []byte) error {
	var tmp uint16
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U16(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u U16) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint16(u))
}

// U32 is an unsigned 32-bit integer
type U32 uint32

// NewU32 creates a new U32 type
func NewU32(u uint32) U32 {
	return U32(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U32) UnmarshalJSON(b []byte) error {
	var tmp uint32
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U32(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u U32) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint32(u))
}

// U64 is an unsigned 64-bit integer
type U64 uint64

// NewU64 creates a new U64 type
func NewU64(u uint64) U64 {
	return U64(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U64) UnmarshalJSON(b []byte) error {
	var tmp uint64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U64(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u U64) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(u))
}

// U128 is an unsigned 128-bit integer, it is represented as a big.Int in Go.
type U128 struct {
	*big.Int
}

// NewU128 creates a new U128 type
func NewU128(i big.Int) U128 {
	return U128{&i}
}

// Decode implements decoding as per the Scale specification
func (i *U128) Decode(decoder Decoder) error {
	bs := make([]byte, 16)
	err := decoder.Read(bs)
	if err != nil {
		return err
	}
	// reverse bytes, scale uses little-endian encoding, big.int's bytes are expected in big-endian
	reverse(bs)

	b, err := UintBytesToBigInt(bs)
	if err != nil {
		return err
	}

	// deal with zero differently to get a nil representation (this is how big.Int deals with 0)
	if b.Sign() == 0 {
		*i = U128{big.NewInt(0)}
		return nil
	}

	*i = U128{b}
	return nil
}

// Encode implements encoding as per the Scale specification
func (i U128) Encode(encoder Encoder) error {
	b, err := BigIntToUintBytes(i.Int, 16)
	if err != nil {
		return err
	}

	// reverse bytes, scale uses little-endian encoding, big.int's bytes are expected in big-endian
	reverse(b)

	return encoder.Write(b)
}

// U256 is an usigned 256-bit integer, it is represented as a big.Int in Go.
type U256 struct {
	*big.Int
}

// NewU256 creates a new U256 type
func NewU256(i big.Int) U256 {
	return U256{&i}
}

// Decode implements decoding as per the Scale specification
func (i *U256) Decode(decoder Decoder) error {
	bs := make([]byte, 32)
	err := decoder.Read(bs)
	if err != nil {
		return err
	}
	// reverse bytes, scale uses little-endian encoding, big.int's bytes are expected in big-endian
	reverse(bs)

	b, err := UintBytesToBigInt(bs)
	if err != nil {
		return err
	}

	// deal with zero differently to get a nil representation (this is how big.Int deals with 0)
	if b.Sign() == 0 {
		*i = U256{big.NewInt(0)}
		return nil
	}

	*i = U256{b}
	return nil
}

// Encode implements encoding as per the Scale specification
func (i U256) Encode(encoder Encoder) error {
	b, err := BigIntToUintBytes(i.Int, 32)
	if err != nil {
		return err
	}

	// reverse bytes, scale uses little-endian encoding, big.int's bytes are expected in big-endian
	reverse(b)

	return encoder.Write(b)
}

// BigIntToUintBytes encodes the given big.Int to a big endian encoded unsigned integer byte slice of the given byte
// length, returning an error if the given big.Int would be bigger than the maximum number the byte slice of the given
// length could hold
func BigIntToUintBytes(i *big.Int, bytelen int) ([]byte, error) {
	if i.Sign() < 0 {
		return nil, fmt.Errorf("cannot encode a negative big.Int into an unsigned integer")
	}

	max := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(bytelen*8)), nil)
	if i.CmpAbs(max) > 0 {
		return nil, fmt.Errorf("cannot encode big.Int to []byte: given big.Int exceeds highest number "+
			"%v for an uint with %v bits", max, bytelen*8)
	}

	res := make([]byte, bytelen)

	bs := i.Bytes()
	copy(res[len(res)-len(bs):], bs)
	return res, nil
}

// UintBytesToBigInt decodes the given byte slice containing a big endian encoded unsigned integer to a big.Int
func UintBytesToBigInt(b []byte) (*big.Int, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("cannot decode an empty byte slice")
	}

	return big.NewInt(0).SetBytes(b), nil
}

// reverse reverses bytes in place (manipulates the underlying array)
func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}