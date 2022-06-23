//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import "encoding/binary"

var (
	MTPPrefix      = []byte{0x01}
	MTPCountPrefix = []byte{0x02}
	ParamsPrefix   = []byte{0x03}
)

func GetMTPKey(address string, id uint64) []byte {
	return append(MTPPrefix, append([]byte(address), GetIDBytes(id)...)...)
}

// GetIDBytes returns the byte representation of the ID
func GetIDBytes(ID uint64) []byte {
	IDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(IDBz, ID)
	return IDBz
}

// GetIDFromBytes returns ID in uint64 format from a byte array
func GetIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
