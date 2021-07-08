package ethbridge

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const TestAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"

func CreateTestHandler(
	t *testing.T, consensusNeeded float64, validatorAmounts []int64,
) (sdk.Context, oracle.Keeper, bank.Keeper, supply.Keeper, auth.AccountKeeper, []sdk.ValAddress, sdk.Handler) {
	ctx, _, oracleKeeper, bankKeeper, supplyKeeper,
		accountKeeper, validatorAddresses, keyEthBridge := keeper.CreateTestKeepers(t, consensusNeeded, validatorAmounts, ModuleName)
	bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	supplyKeeper.SetModuleAccount(ctx, bridgeAccount)

	cdc := keeper.MakeTestCodec()
	// keyEthBridge := sdk.NewKVStoreKey(types.StoreKey)
	bridgeKeeper := NewKeeper(cdc, supplyKeeper, oracleKeeper, keyEthBridge)
	CethReceiverAccount, _ := sdk.AccAddressFromBech32(TestAddress)
	bridgeKeeper.SetCethReceiverAccount(ctx, CethReceiverAccount)
	handler := NewHandler(accountKeeper, bridgeKeeper, cdc)

	return ctx, oracleKeeper, bankKeeper, supplyKeeper, accountKeeper, validatorAddresses, handler
}
