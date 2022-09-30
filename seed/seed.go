package seed

import (
	"math"
	"math/rand"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func FundPool(clpKeeper clpkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	registryKeeper tokenregistrytypes.Keeper,
	ctx sdk.Context,
	denom string,
	address sdk.AccAddress) error {

	pool, err := clpKeeper.GetPool(ctx, denom)
	if err != nil {
		return err
	}

	registry := registryKeeper.GetRegistry(ctx)
	entry, err := registryKeeper.GetEntry(registry, denom)
	if err != nil {
		return err
	}

	native := sdk.NewUint(100).MulUint64(1000000000000000000)
	external := sdk.NewInt(100).MulRaw(int64(math.Pow(10, float64(entry.Decimals))))
	external = pool.SwapPriceNative.MulInt(external).TruncateInt()
	externalUint := sdk.NewUintFromBigInt(external.BigInt())

	msgServer := clpkeeper.NewMsgServerImpl(clpKeeper)

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
		ExternalAssetAmount: externalUint,
	})
	if err != nil {
		return err
	}

	return nil
}

// Seed will generate numAccounts, fund them with enough tokens to add liquidity to the specified pools.
func Seed(clpKeeper clpkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	registryKeeper tokenregistrytypes.Keeper,
	ctx sdk.Context,
	numAccounts int,
	fundNPools int) error {

	// generate numAccounts random addresses, not guaranteeing uniqueness.
	testAddrs := make([]sdk.AccAddress, numAccounts)
	for i := 0; i < numAccounts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	pools := clpKeeper.GetPools(ctx)

	for _, address := range testAddrs {
		for a := 0; a < fundNPools; a++ {
			// select a random pool
			n := rand.Intn(len(pools))
			// fund pool from address
			err := FundPool(clpKeeper, bankKeeper, registryKeeper, ctx, pools[n].ExternalAsset.Symbol, address)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
