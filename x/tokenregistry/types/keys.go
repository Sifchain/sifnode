package types

var WhitelistStorePrefix = []byte{0x01}
var AdminAccountStorePrefix = []byte{0x02}

func GetAdminAccountKey(moduleName string) []byte {
	key := []byte(moduleName)
	return append(AdminAccountStorePrefix, key...)
}
