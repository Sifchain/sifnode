package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	//"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

//nolint:deadcode,unused
var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())

	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpPk2    = ed25519.GenPrivKey().PubKey()
	valOpPk3    = ed25519.GenPrivKey().PubKey()
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address())
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())

	TestAddrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

// test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initPower int64) (
	sdk.Context, Keeper) {
	ctx, keeper := CreateTestInputAdvanced(t, isCheckTx, initPower)
	return ctx, keeper
}

func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, Keeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyClp := sdk.NewKVStoreKey(types.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distribution.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyClp, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[distrAcc.GetAddress().String()] = true

	cdc := MakeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	sk := staking.NewKeeper(cdc, keyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultParams())
	require.Equal(t, staking.DefaultParams(), sk.GetParams(ctx))
	keeper := NewKeeper(cdc, keyClp, bankKeeper, pk.Subspace(types.DefaultParamspace))
	keeper.SetParams(ctx, types.DefaultParams())
	initCoins := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	return ctx, keeper
}

func generateRandomPool(numberOfPools int) []types.Pool {
	var poolList []types.Pool
	tokens := []string{"ETH", "BTC", "EOS", "BCH", "BNB", "USDT", "ADA", "TRX"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfPools; i++ {
		// initialize global pseudo random generator
		externalToken := tokens[rand.Intn(len(tokens))]
		externalAsset := types.NewAsset("ROWAN", "c"+"ROWAN"+externalToken, externalToken)
		pool := types.NewPool(externalAsset, 1, 1, 1)
		poolList = append(poolList, pool)
	}
	return poolList
}

func generateRandomLP(numberOfLp int) []types.LiquidityProvider {
	var lpList []types.LiquidityProvider
	tokens := []string{"ETH", "BTC", "EOS", "BCH", "BNB", "USDT", "ADA", "TRX"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfLp; i++ {
		externalToken := tokens[rand.Intn(len(tokens))]
		asset := types.NewAsset("ROWAN", "c"+"ROWAN"+externalToken, externalToken)
		lp := types.NewLiquidityProvider(asset, 1, "192.0.1.1")
		lpList = append(lpList, lp)
	}
	return lpList
}
