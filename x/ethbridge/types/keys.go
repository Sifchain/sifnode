package types

const (
	// ModuleName is the name of the ethereum bridge module
	ModuleName = "ethbridge"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the ethereum bridge module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the ethereum bridge module
	RouterKey = ModuleName

	// PeggyTokenKey is the key for peggy token list
	PeggyTokenKey = StoreKey + "PeggyToken"

	// ceth symbol
	CethSymbol = "ceth"

	// native symbol
	NativeSymbol = "rowan"
)

var (
	PeggyTokenKeyPrefix       = []byte{0x00}
	CethReceiverAccountPrefix = []byte{0x01}
)
