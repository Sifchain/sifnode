package ethbridge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	ethbridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	oraclekeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
)

const (
	TestAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
)

func CreateTestHandler(t *testing.T, consensusNeeded float64, validatorAmounts []int64) (sdk.Context, ethbridgekeeper.Keeper, bankkeeper.Keeper, oraclekeeper.Keeper, sdk.Handler, []sdk.ValAddress) {

	ctx, keeper, bankKeeper, _, oracleKeeper, _, validatorAddresses := ethbridgekeeper.CreateTestKeepers(t, consensusNeeded, validatorAmounts, "")

	CethReceiverAccount, _ := sdk.AccAddressFromBech32(TestAddress)
	keeper.SetCethReceiverAccount(ctx, CethReceiverAccount)
	handler := NewHandler(keeper)

	return ctx, keeper, bankKeeper, oracleKeeper, handler, validatorAddresses
}
