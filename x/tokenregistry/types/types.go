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
	case "permission_clp":
		return Permission_CLP
	case "permission_ibc_export":
		return Permission_IBCEXPORT
	case "permission_ibc_import":
		return Permission_IBCIMPORT
	default:
		return Permission_UNSPECIFIED
	}
}

func StringCompare(a, b string) bool {
	return a == b
}
