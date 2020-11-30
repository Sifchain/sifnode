package common

import (
	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/tools/sifgen/common/types"
)

// Aliases
type (
	Keys       types.Keys
	NodeConfig types.NodeConfig
	CLIConfig  types.CLIConfig
	Genesis    types.Genesis
)

var (
	DefaultNodeHome = app.DefaultNodeHome
	DefaultCLIHome  = app.DefaultCLIHome
	StakeTokenDenom = types.StakeTokenDenom
)

var (
	ToFund = types.ToFund
	ToBond = types.ToBond
)
