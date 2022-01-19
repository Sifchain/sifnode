package types

const (
	// ModuleName is the name of the whitelist module
	ModuleName = "tokenregistry"

	// StoreKey is the string store representation
	StoreKey = ModuleName + "-2"

	// QuerierRoute is the querier route
	QuerierRoute = ModuleName

	// RouterKey is the msg router key
	RouterKey = ModuleName
)

func GetPermissionFromString(s string) Permission {
	switch s {
	case "PERMISSION_CLP":
		return Permission_PERMISSION_CLP
	case "PERMISSION_IBCEXPORT":
		return Permission_PERMISSION_IBCEXPORT
	case "PERMISSION_IBCIMPORT":
		return Permission_PERMISSION_IBCIMPORT
	default:
		return Permission_PERMISSION_UNSPECIFIED
	}
}
