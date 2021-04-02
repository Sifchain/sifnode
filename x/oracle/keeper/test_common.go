package keeper

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramkeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	TestID                     = "oracleID"
	AlternateTestID            = "altOracleID"
	TestString                 = "{value: 5}"
	AlternateTestString        = "{value: 7}"
	AnotherAlternateTestString = "{value: 9}"
)

// CreateTestKeepers greates an Mock App, OracleKeeper, bankKeeper and ValidatorAddresses to be used for test input
func CreateTestKeepers(t *testing.T, consensusNeeded float64, validatorAmounts []int64, extraMaccPerm string) (
	sdk.Context, Keeper, bankKeeper.Keeper, authkeeper.AccountKeeper, []sdk.ValAddress, sdk.StoreKey) {
	PKs := CreateTestPubKeys(500)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramtypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramtypes.TStoreKey)
	keySupply := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyEthBridge := sdk.NewKVStoreKey(extraMaccPerm)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyEthBridge, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "foochainid"}, false, nil)
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &tmproto.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	ctx = ctx.WithLogger(log.NewNopLogger())
	cdc := MakeTestCodec()

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	paramsKeeper := paramkeeper.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		paramsKeeper.Subspace(auth.DefaultParamspace),
		authtypes.ProtoBaseAccount, // prototype
	)

	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	}

	if extraMaccPerm != "" {
		maccPerms[extraMaccPerm] = []string{authtypes.Burner, authtypes.Minter}
	}

	initTokens := sdk.TokensFromConsensusPower(10000)
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(100))))

	stakingKeeper := stakingkeeper.NewKeeper(cdc, keyStaking, supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace))
	stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())
	oracleKeeper := NewKeeper(cdc, keyOracle, stakingKeeper, consensusNeeded)

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	authtypes.SetModuleAccount(ctx, feeCollectorAcc)
	authtypes.SetModuleAccount(ctx, bondPool)
	authtypes.SetModuleAccount(ctx, notBondedPool)

	// Setup validators
	valAddrs := make([]sdk.ValAddress, len(validatorAmounts))
	for i, amount := range validatorAmounts {
		valPubKey := PKs[i]
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
		valAddrs[i] = valAddr
		valTokens := sdk.TokensFromConsensusPower(amount)
		// test how the validator is set from a purely unbonbed pool
		validator := stakingtypes.NewValidator(valAddr, valPubKey, stakingtypes.Description{})
		validator, _ = validator.AddTokensFromDel(valTokens)
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	oracleKeeper.SetOracleWhiteList(ctx, valAddrs)
	return ctx, oracleKeeper, bankKeeper, accountKeeper, valAddrs, keyEthBridge
}

// nolint: unparam
func CreateTestAddrs(numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	var addresses []sdk.AccAddress
	var valAddresses []sdk.ValAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		address := stakingkeeper.TestAddr(buffer.String(), bech)
		valAddress := sdk.ValAddress(address)
		addresses = append(addresses, address)
		valAddresses = append(valAddresses, valAddress)
		buffer.Reset()
	}
	return addresses, valAddresses
}

// nolint: unparam
func CreateTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString(
			"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF",
		) //base pubkey string
		buffer.WriteString(numString) //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKey
	copy(pkEd[:], pkBytes)
	return pkEd
}

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	// Register AppAccount
	cdc.RegisterInterface((*authtypes.AccountI)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/staking/BaseAccount", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	staking.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	return cdc
}
