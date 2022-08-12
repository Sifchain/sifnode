//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import "encoding/binary"

var (
	MTPPrefix          = []byte{0x01}
	MTPCountPrefix     = []byte{0x02}
	ParamsPrefix       = []byte{0x03}
	OpenMTPCountPrefix = []byte{0x04}
	WhitelistPrefix    = []byte{0x05}
)

func GetMTPKey(address string, id uint64) []byte {
	return append(MTPPrefix, append([]byte(address), GetUint64Bytes(id)...)...)
}

func GetMTPPrefixForAddress(address string) []byte {
	return append(MTPPrefix, []byte(address)...)
}

// GetUint64Bytes returns the byte representation of the ID
func GetUint64Bytes(ID uint64) []byte {
	IDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(IDBz, ID)
	return IDBz
}

// GetUint64FromBytes returns ID in uint64 format from a byte array
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

func GetWhitelistKey(address string) []byte {
	return append(WhitelistPrefix, []byte(address)...)
}
