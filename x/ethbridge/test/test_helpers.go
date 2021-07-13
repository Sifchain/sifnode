package test

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"time"

	"strconv"
	"testing"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/keeper"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oraclekeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	oracleTypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	TestID                         = "oracleID"
	AlternateTestID                = "altOracleID"
	TestString                     = "{value: 5}"
	AlternateTestString            = "{value: 7}"
	AnotherAlternateTestString     = "{value: 9}"
	TestNativeTokenReceiverAddress = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv" //nolint
	NetworkDescriptor              = 1
	NativeToken                    = "ceth"
	NativeTokenGas                 = 1
	MinimumCost                    = 1
)

// CreateTestKeepers greates an Mock App, OracleKeeper, bankKeeper and ValidatorAddresses to be used for test input
func CreateTestKeepers(t *testing.T, consensusNeeded float64, validatorAmounts []int64, extraMaccPerm string) (sdk.Context, keeper.Keeper, bankkeeper.Keeper, authkeeper.AccountKeeper, oraclekeeper.Keeper,
	simappparams.EncodingConfig, oracleTypes.ValidatorWhiteList, []sdk.ValAddress) {

	PKs := CreateTestPubKeys(500)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	// TODO: staking.TStoreKey removed in favor of?
	tkeyStaking := sdk.NewTransientStoreKey("transient_staking")
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracleTypes.StoreKey)
	keyEthBridge := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
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
	encCfg := MakeTestEncodingConfig()

	bridgeAccount := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Burner, authtypes.Minter)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Burner, authtypes.Minter},
	}

	if extraMaccPerm != "" {
		maccPerms[extraMaccPerm] = []string{authtypes.Burner, authtypes.Minter}
	}

	paramsKeeper := paramskeeper.NewKeeper(encCfg.Marshaler, encCfg.Amino, keyParams, tkeyParams)

	//accountKeeper gets maccParams in 0.40, module accounts moved from supplykeeper to authkeeper
	accountKeeper := authkeeper.NewAccountKeeper(
		encCfg.Marshaler, // amino codec
		keyAcc,           // target store
		paramsKeeper.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, // prototype,
		maccPerms,
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		encCfg.Marshaler,
		keyBank,
		accountKeeper,
		paramsKeeper.Subspace(banktypes.ModuleName),
		blacklistedAddrs,
	)

	initTokens := sdk.TokensFromConsensusPower(10000)
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(100))))

	bankKeeper.SetSupply(ctx, banktypes.NewSupply(totalSupply))

	stakingKeeper := stakingkeeper.NewKeeper(encCfg.Marshaler, keyStaking, accountKeeper, bankKeeper, paramsKeeper.Subspace(stakingtypes.ModuleName))
	stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())
	oracleKeeper := oraclekeeper.NewKeeper(encCfg.Marshaler, keyOracle, stakingKeeper, consensusNeeded)

	// set module accounts
	err = bankKeeper.AddCoins(ctx, notBondedPool.GetAddress(), totalSupply)
	require.NoError(t, err)

	accountKeeper.SetModuleAccount(ctx, bridgeAccount)
	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)

	ethbridgeKeeper := keeper.NewKeeper(encCfg.Marshaler, bankKeeper, oracleKeeper, accountKeeper, keyEthBridge)

	NativeTokenReceiverAccount, _ := sdk.AccAddressFromBech32(TestNativeTokenReceiverAddress)
	ethbridgeKeeper.SetNativeTokenReceiverAccount(ctx, NativeTokenReceiverAccount)

	// Setup validators
	valAddrsInOrder := make([]sdk.ValAddress, len(validatorAmounts))
	valAddrs := make(map[string]uint32)
	for i, amount := range validatorAmounts {
		valPubKey := PKs[i]
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
		valAddrsInOrder[i] = valAddr
		valAddrs[valAddr.String()] = uint32(amount)
		valTokens := sdk.TokensFromConsensusPower(amount)
		// test how the validator is set from a purely unbonbed pool
		validator, err := stakingtypes.NewValidator(valAddr, valPubKey, stakingtypes.Description{})
		require.NoError(t, err)

		validator, _ = validator.AddTokensFromDel(valTokens)
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		_, err = stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
		if err != nil {
			panic("Failed to apply validator set updates")
		}
	}

	networkIdentity := oracleTypes.NewNetworkIdentity(NetworkDescriptor)

	oracleKeeper.SetNativeToken(ctx, networkIdentity, NativeToken,
		sdk.NewInt(NativeTokenGas), sdk.NewInt(MinimumCost), sdk.NewInt(MinimumCost))
	whitelist := oracleTypes.ValidatorWhiteList{WhiteList: valAddrs}
	oracleKeeper.SetOracleWhiteList(ctx, networkIdentity, whitelist)

	return ctx, ethbridgeKeeper, bankKeeper, accountKeeper, oracleKeeper, encCfg, whitelist, valAddrsInOrder
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
		address, _ := sdk.AccAddressFromHex(buffer.String())
		valAddress := sdk.ValAddress(address)
		addresses = append(addresses, address)
		valAddresses = append(valAddresses, valAddress)
		buffer.Reset()
	}
	return addresses, valAddresses
}

// nolint: unparam
func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
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

func NewPubKey(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}

	//res, err = crypto.PubKeyFromBytes(pkBytes)
	return &ed25519.PubKey{
		Key: pkBytes,
	}
}

// create a codec used only for testing
func MakeTestEncodingConfig() simappparams.EncodingConfig {
	return app.MakeTestEncodingConfig()
}

const (
	AddressKey1 = "A58856F0FD53BF058B4909A21AEC019107BA6"
)

//// returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*app.SifchainApp, sdk.Context) {
	sifapp := app.Setup(isCheckTx)
	ctx := sifapp.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	sifapp.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)

	_ = app.AddTestAddrs(sifapp, ctx, 6, initTokens)

	return sifapp, ctx
}

func CreateTestAppEthBridge(isCheckTx bool) (sdk.Context, keeper.Keeper) {
	sifapp, ctx := CreateTestApp(isCheckTx)
	return ctx, sifapp.EthbridgeKeeper
}

func GenerateRandomTokens(numberOfTokens int) []string {
	var tokenList []string
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx", "cacoin", "cbcoin", "ccoin", "cdcoin"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfTokens; i++ {
		// initialize global pseudo random generator
		randToken := tokens[rand.Intn(len(tokens))]

		tokenList = append(tokenList, randToken)
	}
	return tokenList
}

func GenerateAddress(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString(strconv.Itoa(100))
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr := buffer.String()
	res, err := sdk.AccAddressFromHex(addr)

	if err != nil {
		panic(err)
	}

	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}
	return res
}
