package types

const (
	// ModuleName is the name of the whitelist module
	ModuleName = "admin"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route
	QuerierRoute = ModuleName

	// RouterKey is the msg router key
	RouterKey = ModuleName
)

func StringCompare(a, b string) bool {
	return a == b
}
