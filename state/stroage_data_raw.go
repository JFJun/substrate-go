package state

import (
	"fmt"
	"github.com/JFJun/substrate-go/scale"
	"io"
)

type StorageDataRaw []byte

// NewStorageDataRaw creates a new StorageDataRaw type
func NewStorageDataRaw(b []byte) StorageDataRaw {
	return StorageDataRaw(b)
}

// Encode implements encoding for StorageDataRaw, which just unwraps the bytes of StorageDataRaw
func (s StorageDataRaw) Encode(encoder scale.Encoder) error {
	return encoder.Write(s)
}

// Decode implements decoding for StorageDataRaw, which just reads all the remaining bytes into StorageDataRaw
func (s *StorageDataRaw) Decode(decoder scale.Decoder) error {
	for i := 0; true; i++ {
		b, err := decoder.ReadOneByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		*s = append((*s)[:i], b)
	}
	return nil
}

// Hex returns a hex string representation of the value
func (s StorageDataRaw) Hex() string {
	return fmt.Sprintf("%#x", s)
}
