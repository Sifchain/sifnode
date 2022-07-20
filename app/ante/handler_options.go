package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// HandlerOptions defines the list of module keepers required to run the Sifnode
// AnteHandler decorators.
type HandlerOptions struct {
	AccountKeeper   ante.AccountKeeper
	BankKeeper      authtypes.BankKeeper
	FeegrantKeeper  ante.FeegrantKeeper
	StakingKeeper   stakingkeeper.Keeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
}
