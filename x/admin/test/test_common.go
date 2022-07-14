package test

import (
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
)

func GetAdmins(address string) []*admintypes.AdminAccount {
	return []*admintypes.AdminAccount{
		{
			AdminType:    admintypes.AdminType_ADMIN,
			AdminAddress: address,
		},
		{
			AdminType:    admintypes.AdminType_PMTPREWARDS,
			AdminAddress: address,
		},
		{
			AdminType:    admintypes.AdminType_CLPDEX,
			AdminAddress: address,
		},
		{
			AdminType:    admintypes.AdminType_TOKENREGISTRY,
			AdminAddress: address,
		},
		{
			AdminType:    admintypes.AdminType_ETHBRIDGE,
			AdminAddress: address,
		},
	}
}
