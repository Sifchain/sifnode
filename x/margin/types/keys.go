package types

var MTPPrefix = []byte{0x01}

func GetMTPKey(asset, address string) []byte {
	return append(MTPPrefix, append([]byte(asset), []byte(address)...)...)
}
