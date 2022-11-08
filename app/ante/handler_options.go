package ante

import (
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// HandlerOptions defines the list of module keepers required to run the Sifnode
// AnteHandler decorators.
type HandlerOptions struct {
	AdminKeeper     adminkeeper.Keeper
	AccountKeeper   ante.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	FeegrantKeeper  ante.FeegrantKeeper
	StakingKeeper   stakingkeeper.Keeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
}
