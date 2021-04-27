package common

import (
	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/tools/sifgen/common/types"
)

// Aliases
type (
	Keys       types.Keys
	NodeConfig types.NodeConfig
	Genesis    types.Genesis
)

var (
	DefaultNodeHome = app.DefaultNodeHome
	StakeTokenDenom = types.StakeTokenDenom

	P2PPort             = 26656
	MaxNumInboundPeers  = 1000
	MaxNumOutboundPeers = 1000
	AllowDuplicateIP    = true
)

var (
	MinCLPCreatePoolThreshold = "100"
)
