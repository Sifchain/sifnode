package seed

import (
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func FundPool(clpKeeper clpkeeper.Keeper, bankKeeper bankkeeper.Keeper, ctx sdk.Context, denom string, addresses []sdk.AccAddress) error {
	pool, err := clpKeeper.GetPool(ctx, denom)
	if err != nil {
		return err
	}

	native := pool.NativeAssetBalance
	external := pool.ExternalAssetBalance

	msgServer := clpkeeper.NewMsgServerImpl(clpKeeper)

	for _, address := range addresses {
		err = simapp.FundAccount(bankKeeper, ctx, address, sdk.NewCoins(
			sdk.NewCoin("rowan", sdk.NewIntFromBigInt(native.BigInt())),
			sdk.NewCoin(denom, sdk.NewIntFromBigInt(external.BigInt())),
		))
		if err != nil {
			return err
		}

		_, err = msgServer.AddLiquidity(sdk.WrapSDKContext(ctx), &clptypes.MsgAddLiquidity{
			Signer:              address.String(),
			ExternalAsset:       &clptypes.Asset{Symbol: denom},
			NativeAssetAmount:   native,
			ExternalAssetAmount: external,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Seed will generate numAccounts, fund them with enough tokens to add liquidity to the specified pools.
func Seed(clpKeeper clpkeeper.Keeper, bankKeeper bankkeeper.Keeper, ctx sdk.Context, numAccounts int, denoms []string) error {
	testAddrs := make([]sdk.AccAddress, numAccounts)
	for i := 0; i < numAccounts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	for _, denom := range denoms {
		err := FundPool(clpKeeper, bankKeeper, ctx, denom, testAddrs)
		if err != nil {
			return err
		}
	}

	return nil
}
