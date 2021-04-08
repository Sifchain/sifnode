package ethbridge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keeperLib "github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/Sifchain/sifnode/x/oracle"
)

const (
	TestAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
)

func CreateTestHandler(
	t *testing.T, consensusNeeded float64, validatorAmounts []int64,
) (sdk.Context, oracle.Keeper, bank.Keeper, supply.Keeper, auth.AccountKeeper, []sdk.ValAddress, sdk.Handler) {
	ctx, oracleKeeper, bankKeeper, supplyKeeper,
		accountKeeper, validatorAddresses, keyEthBridge := oracle.CreateTestKeepers(t, consensusNeeded, validatorAmounts, ModuleName)
	bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	supplyKeeper.SetModuleAccount(ctx, bridgeAccount)

	cdc := keeperLib.MakeTestCodec()
	// keyEthBridge := sdk.NewKVStoreKey(types.StoreKey)
	bridgeKeeper := NewKeeper(cdc, supplyKeeper, oracleKeeper, keyEthBridge)
	CethReceiverAccount, _ := sdk.AccAddressFromBech32(TestAddress)
	bridgeKeeper.SetCethReceiverAccount(ctx, CethReceiverAccount)
	handler := NewHandler(bridgeKeeper)

	return ctx, oracleKeeper, bankKeeper, supplyKeeper, accountKeeper, validatorAddresses, handler
}
