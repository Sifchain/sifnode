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
)

var (
	PeggyTokenKeyPrefix                = []byte{0x00}
	CrossChainFeeReceiverAccountPrefix = []byte{0x01}
	BlacklistPrefix                    = []byte{0x02}
	PausePrefix                        = []byte{0x03}
	GlobalNoncePrefix                  = []byte{0x04}
	EthereumLockBurnSequencePrefix     = []byte{0x05}
	GlobalNonceToBlockNumberPrefix     = []byte{0x06}
	FirstLockDoublePegPrefix           = []byte{0x07}
)
