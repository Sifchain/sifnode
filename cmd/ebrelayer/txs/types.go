package txs

import (
	"math/big"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
)

// OracleClaim contains data required to make an OracleClaim
type OracleClaim struct {
	ProphecyID *big.Int
	Message    [32]byte
	Signature  []byte
}

// ProphecyClaim contains data required to make an ProphecyClaim
type ProphecyClaim struct {
	CosmosSender     []byte
	Symbol           string
	Amount           *big.Int
	EthereumReceiver common.Address
	ClaimType        types.Event
}
