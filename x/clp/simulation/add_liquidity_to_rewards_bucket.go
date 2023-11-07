package simulation

import (
	"math/rand"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddLiquidityToRewardsBucket(
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddLiquidityToRewardsBucketRequest{
			Signer: simAccount.Address.String(),
		}

		// TODO: Handling the AddLiquidityToRewardsBucket simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AddLiquidityToRewardsBucket simulation not implemented"), nil, nil
	}
}
