package types

import "fmt"

var AdminAccountStorePrefix = []byte{0x01}

func GetAdminAccountKey(adminAccount AdminAccount) []byte {
	key := []byte(fmt.Sprintf("%s_%s", adminAccount.AdminType.String(), adminAccount.AdminAddress))
	return append(AdminAccountStorePrefix, key...)
}
