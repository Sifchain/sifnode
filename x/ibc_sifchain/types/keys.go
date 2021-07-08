package types

import "fmt"

const StoreKey = "sifchain_ibc"

var (
	WhiteListPrefix = []byte{0x0001} // key for storing WhiteList in state
)

func GetWhiteListKey() []byte {
	key := []byte(fmt.Sprintf("%s", "Sifchain_denom_whitelist"))
	return append(WhiteListPrefix, key...)
}
