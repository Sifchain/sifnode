package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// SimApp testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func setup(withGenesis bool, invCheckPeriod uint, blacklist []sdk.AccAddress) (*SifchainApp, GenesisState) {
	db := dbm.NewMemDB()
	encCdc := MakeTestEncodingConfig()
	app := NewSifAppWithBlacklist(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encCdc, EmptyAppOptions{}, blacklist)
	if withGenesis {
		return app, NewDefaultGenesisState(encCdc.Marshaler)
	}
	return app, GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *SifchainApp {
	return SetupWithBlacklist(isCheckTx, []sdk.AccAddress{})
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func SetupWithBlacklist(isCheckTx bool, blacklist []sdk.AccAddress) *SifchainApp {
	app, genesisState := setup(!isCheckTx, 5, blacklist)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}
	return app
}

func SetupFromGenesis(isCheckTx bool, genesisTransformer func(*SifchainApp, GenesisState) GenesisState) *SifchainApp {
	app, genesisState := setup(!isCheckTx, 5, []sdk.AccAddress{})

	genesisState = genesisTransformer(app, genesisState)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}
	return app
}

// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") // base pubkey string
		buffer.WriteString(numString)                                                       // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	return valAddrs
}

// NewPubKeyFromHex returns a PubKey from a hex string.
func NewPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		panic(errors.Wrap(errors.ErrInvalidPubKey, "invalid pubkey size"))
	}
	return &ed25519.PubKey{Key: pkBytes}
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// AddTestAddrs constructs and returns accNum amount of accounts with initial balance of accAmt in random order
func AddTestAddrs(app *SifchainApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, CreateRandomAccounts)
}

func addTestAddrs(app *SifchainApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)
	moduleCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(accNum+1))))
	_ = app.BankKeeper.MintCoins(ctx, "ethbridge", moduleCoins)
	accCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	for _, addr := range testAddrs {
		err := app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "ethbridge", addr, accCoins)
		if err != nil {
			panic(err)
		}
	}
	return testAddrs
}

// AddCoinsToAccount initializes an account with balance based on initCoins.
func AddCoinsToAccount(moduleName string, bankKeeper keeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) error {
	err := bankKeeper.MintCoins(ctx, moduleName, initCoins)
	if err != nil {
		return err
	}
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, addr, initCoins)
	return err
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// CreateRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func CreateRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}
	return testAddrs
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}
	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}
	return res, nil
}
