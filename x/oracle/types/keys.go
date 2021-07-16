package types

const (
	// ModuleName is the name of the oracle module
	ModuleName = "oracle"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the oracle module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the oracle module
	RouterKey = ModuleName
)

var (
	WhiteListValidatorPrefix = []byte{0x00}
	AdminAccountPrefix       = []byte{0x01}
	ProphecyPrefix           = []byte{0x02}
	CrossChainFeePrefix      = []byte{0x03}
)
