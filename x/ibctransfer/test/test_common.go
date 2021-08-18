package test

// import (
// 	"testing"

// 	sifapp "github.com/Sifchain/sifnode/app"
// 	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
// 	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	"github.com/cosmos/cosmos-sdk/x/bank/types"
// 	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
// 	"github.com/cosmos/cosmos-sdk/x/ibc/testing/mock"
// 	"github.com/stretchr/testify/require"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
// 	tmtypes "github.com/tendermint/tendermint/types"
// )

// type GenesisState map[string]json.RawMessage

// func NewDefaultGenesisState(cdc codec.JSONMarshaler) GenesisState {
// 	return ModuleBasics.DefaultGenesis(cdc)
// }

// func setup(withGenesis bool, invCheckPeriod uint) (*sifapp.SifchainApp, GenesisState) {
// 	db := dbm.NewMemDB()
// 	encCdc := MakeTestEncodingConfig()
// 	app := NewSifApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encCdc, EmptyAppOptions{})
// 	if withGenesis {
// 		return app, NewDefaultGenesisState(encCdc.Marshaler)
// 	}
// 	return app, GenesisState{}
// }

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
// func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *SimApp {
// 	app, genesisState := setup(true, 5)
// 	// set genesis accounts
// 	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
// 	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

// 	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
// 	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

// 	bondAmt := sdk.NewInt(1000000)

// 	for _, val := range valSet.Validators {
// 		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
// 		require.NoError(t, err)
// 		pkAny, err := codectypes.NewAnyWithValue(pk)
// 		require.NoError(t, err)
// 		validator := stakingtypes.Validator{
// 			OperatorAddress:   sdk.ValAddress(val.Address).String(),
// 			ConsensusPubkey:   pkAny,
// 			Jailed:            false,
// 			Status:            stakingtypes.Bonded,
// 			Tokens:            bondAmt,
// 			DelegatorShares:   sdk.OneDec(),
// 			Description:       stakingtypes.Description{},
// 			UnbondingHeight:   int64(0),
// 			UnbondingTime:     time.Unix(0, 0).UTC(),
// 			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
// 			MinSelfDelegation: sdk.ZeroInt(),
// 		}
// 		validators = append(validators, validator)
// 		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

// 	}

// 	// set validators and delegations
// 	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
// 	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

// 	totalSupply := sdk.NewCoins()
// 	for _, b := range balances {
// 		// add genesis acc tokens and delegated tokens to total supply
// 		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))...)
// 	}

// 	// update total supply
// 	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
// 	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

// 	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
// 	require.NoError(t, err)

// 	// init chain will set the validator set and initialize the genesis accounts
// 	app.InitChain(
// 		abci.RequestInitChain{
// 			Validators:      []abci.ValidatorUpdate{},
// 			ConsensusParams: DefaultConsensusParams,
// 			AppStateBytes:   stateBytes,
// 		},
// 	)

// 	// commit genesis changes
// 	app.Commit()
// 	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
// 		Height:             app.LastBlockHeight() + 1,
// 		AppHash:            app.LastCommitID().Hash,
// 		ValidatorsHash:     valSet.Hash(),
// 		NextValidatorsHash: valSet.Hash(),
// 	}})

// 	return app
// }
