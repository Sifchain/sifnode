# Sifchain - Native Integration Test Documentation

##### Dependencies

0. `git clone git@github.com:Sifchain/sifnode.git`
1. `cd sifnode`
2. `git checkout feature/ptmp-integration-setup-env`

#### What is native integration testing

Cosmos SDK simulator app [simapp](https://docs.cosmos.network/master/building-modules/simulator.html) allows to initiate and run a simulated cosmos chain instance within a testing environment in native go language. This is enabled by the cosmos SDK libraries [testutil/network](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/testutil/network) which both provide all the components to create a new chain network and run the chain instance within a test suite.

The first time the simapp chain is initiated, the whole process is cached allowing instantanious runs the next times they are played and therefore making the whole integration testing experience more conveniant.

The following documentation will explain how to create an integration test environment for the PMTP module.

#### Setup

1. Change the working directory to the following path;

```bash
cd x/pmtp/client/testutil
```

2. Try running the existing integration test suite available;

```bash
go test $(go list ./... | grep -v /vendor/) -tags=integration -v --run TestIntegrationTestSuite
```

Result:

```
=== RUN   TestIntegrationTestSuite
    suite.go:32: setting up integration test suite
    network.go:173: acquiring test network lock
    network.go:178: created temporary directory: /tmp/TestIntegrationTestSuite1190615240/001/chain-IKdhtE2426805842
    network.go:187: preparing test network...
    network.go:381: starting test network...
[...]
    network.go:386: started test network
=== RUN   TestIntegrationTestSuite/TestCLPsExists
=== RUN   TestIntegrationTestSuite/TestRowanBalanceExists
=== CONT  TestIntegrationTestSuite
    suite.go:71: tearing down integration test suite
    network.go:473: cleaning up test network...
    network.go:496: finished cleaning up test network
    network.go:470: released test network lock
--- PASS: TestIntegrationTestSuite (16.46s)
    --- PASS: TestIntegrationTestSuite/TestCLPsExists (0.00s)
    --- PASS: TestIntegrationTestSuite/TestRowanBalanceExists (0.00s)
PASS
ok      github.com/Sifchain/sifnode/x/pmtp/client/testutil      (cached)
```

3. You can find all the current integration test cases defined within the `suite.go` file;

```bash
cat suite.go
```

#### Create a new integration test suite

From now on we are going to walk through each step in order to create a new test suite for integration testing. We also want the test environment to define initial genesis states so we can run tests against pre-defined conditions.

1. Let's create a new test suite file named `cli_test.go` with the following content;

```go
//go:build integration
// +build integration

package testutil

import (
	"os"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestIntegrationTestSuite(t *testing.T) { // name the test suite main function
	cfg := network.DefaultConfig() // retrieve the default testutil/network library config object
	cfg.NumValidators = 1 // limit the number of validators to 1 in our test environment
	encConfig := sifapp.MakeTestEncodingConfig() // retrieve the codecs
	cfg.InterfaceRegistry = encConfig.InterfaceRegistry // set the interface registry
	cfg.Codec = encConfig.Marshaler // set the codec to use
	cfg.TxConfig = encConfig.TxConfig // set the transaction config object
	cfg.AppConstructor = func(val network.Validator) servertypes.Application { // app construction function to create a new chain app based on sifchain app
		return sifapp.NewSifApp(
			log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
			dbm.NewMemDB(),
			nil,
			true,
			make(map[int64]bool),
			val.Dir,
			0,
			encConfig,
			sifapp.EmptyAppOptions{},
			baseapp.SetTrace(true), // add more visibility to the errors thrown during testing
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices), // override min gas price setting
		)
	}

	suite.Run(t, NewIntegrationTestSuite(cfg)) // run the test suite once this test file gets called
}
```

You can refer to the comment above to understand what the code does.

2. Then let's create a `cli_helpers.go` file to define all the necessary sifchain commands we will be using during our tests;

```go
package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	clpcli "github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
)

func QueryBalancesExec(clientCtx client.Context, address fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) { // run query balances command to retrieve all the balances for a particular sifchain address
	args := []string{address.String(), fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), args)
}

func QueryClpPoolsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) { // run query clp pools command to retrieve all the existing liquidity pools
	args := []string{fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdPools(""), args)
}
```

3. Then let's create a `genesis_helpers.go` file to define all the function that override the genesis states for any given module we have to run test against;

```go
package testutil

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// this function creates a subset of the genesis state for the bank module to override the default assigned balances and tokens to the validator address
func GetBankGenesisState(cfg network.Config, address string, nativeAmount sdk.Int, externalAmount sdk.Int) ([]byte, error) {
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("catk", externalAmount),
				sdk.NewCoin("cbtk", externalAmount),
				sdk.NewCoin("cdash", externalAmount),
				sdk.NewCoin("ceth", externalAmount),
				sdk.NewCoin("clink", externalAmount),
				sdk.NewCoin("rowan", nativeAmount),
			},
		},
	}
	gs := banktypes.DefaultGenesisState()
	gs.Balances = append(gs.Balances, balances...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

// this function creates a subset of the genesis state for the tokenregistry module to whitelist all the new tokens that have been added through the bank module
func GetTokenRegistryGenesisState(cfg network.Config, address string) ([]byte, error) {
	gs := &tokenregistrytypes.GenesisState{
		AdminAccount: address,
		Registry: &tokenregistrytypes.Registry{
			Entries: []*tokenregistrytypes.RegistryEntry{
				{Denom: "node0token", BaseDenom: "node0token", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "catk", BaseDenom: "catk", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "cbtk", BaseDenom: "cbtk", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "cdash", BaseDenom: "cdash", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "ceth", BaseDenom: "ceth", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "clink", BaseDenom: "clink", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			},
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

// this function creates a subset of the genesis state for the CLP module to create new liquidy pools for a couple of tokens pairs.
func GetClpGenesisState(cfg network.Config, pool1Amount sdk.Uint, pool2Amount sdk.Uint) ([]byte, error) {
	pools := []*clptypes.Pool{
		{
			ExternalAsset:        &clptypes.Asset{Symbol: "cdash"},
			NativeAssetBalance:   pool1Amount,
			ExternalAssetBalance: pool1Amount,
			PoolUnits:            pool1Amount,
		},
		{
			ExternalAsset:        &clptypes.Asset{Symbol: "ceth"},
			NativeAssetBalance:   pool2Amount,
			ExternalAssetBalance: pool2Amount,
			PoolUnits:            pool2Amount,
		},
	}
	gs := clptypes.DefaultGenesisState()
	gs.PoolList = append(gs.PoolList, pools...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}
```

4. Then let's create a `suite.go` file where all our test suite test cases will be defined;

```go
//go:build integration
// +build integration

package testutil

import (
	sifapp "github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
)

// our test suite structure that includes additional fields such as mnemonic and address as they are later used within the test cases
type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	mnemonic       string
	address        string
	nativeAmount   types.Int
	externalAmount types.Int
}

// main function to initiate our test suite using the config object as it is passed by the `cli_test.go` main function
func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

// we setup the test suite here with a consistent mnemonic and relevant sif address as they are used to allocate additional resource to this address (validator), the chain gets also initiated and started here
func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	sifapp.SetConfig(false)

	s.mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
	s.address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	s.nativeAmount, _ = types.NewIntFromString("999999000000000000000000000")
	s.externalAmount, _ = types.NewIntFromString("500000000000000000000000")

	s.cfg.Mnemonics = []string{s.mnemonic}
	s.cfg.StakingTokens = s.nativeAmount

    // gets the bank genesis state and assign it to our new chain genesis state
	bz, err := GetBankGenesisState(s.cfg, s.address, s.nativeAmount, s.externalAmount)
	s.Require().NoError(err)
	s.cfg.GenesisState["bank"] = bz

    // gets the tokenregistry genesis state and assign it to our new chain genesis state
	bz, err = GetTokenRegistryGenesisState(s.cfg, s.address)
	s.Require().NoError(err)
	s.cfg.GenesisState["tokenregistry"] = bz

    // gets the clp genesis state and assign it to our new chain genesis state
	bz, err = GetClpGenesisState(s.cfg, types.NewUint(3000000000000000000), types.NewUint(2000000000000000000))
	s.Require().NoError(err)
	s.cfg.GenesisState["clp"] = bz

    // init and start the chain
	s.network = network.New(s.T(), s.cfg)

    // we wait until the first block gets written
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// this function is called once all the test cases have been completed, it stops and clear all the data related with the test chain
func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	s.network.Cleanup()
}

// this is our first test case where we want to make sure our validator address who shares the same mnemonic and sif address as the one defined within SetupSuite, holds the relevant amount of balance in the native token
func (s *IntegrationTestSuite) TestRowanBalanceExists() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	var genesisState banktypes.GenesisState
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(s.cfg.GenesisState["bank"], &genesisState))
	s.Require().Equal(genesisState.Balances[0].Address, s.address)
	s.Require().Equal(genesisState.Balances[0].Coins[5], types.NewCoin("rowan", s.nativeAmount))

	out, err := QueryBalancesExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var balancesRes banktypes.QueryAllBalancesResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &balancesRes), out.String())

	s.Require().Contains(balancesRes.Balances, types.NewCoin("rowan", s.nativeAmount))
}

// in our second test case we want to make sure that the initial CLPs in SetupSuite are existings as we call the relevant clp pools command with the relevant pools created
func (s *IntegrationTestSuite) TestCLPsExists() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryClpPoolsExec(clientCtx)
	s.Require().NoError(err)

	var poolsRes clptypes.PoolsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolsRes), out.String())

	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:        &clptypes.Asset{Symbol: "cdash"},
			NativeAssetBalance:   types.NewUint(3000000000000000000),
			ExternalAssetBalance: types.NewUint(3000000000000000000),
			PoolUnits:            types.NewUint(3000000000000000000),
		},
	)
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:        &clptypes.Asset{Symbol: "ceth"},
			NativeAssetBalance:   types.NewUint(2000000000000000000),
			ExternalAssetBalance: types.NewUint(2000000000000000000),
			PoolUnits:            types.NewUint(2000000000000000000),
		},
	)
}
```

5. You can add your own test case by writing a new function within the `suite.go` file;

```go
func (s *IntegrationTestSuite) TestMyNewTestCase() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

    // test code goes here
}
```

6. Try running the integration test suite with the newly added test case;

```bash
go test $(go list ./... | grep -v /vendor/) -tags=integration -v --run TestIntegrationTestSuite
```

Result:

```
=== RUN   TestIntegrationTestSuite
    suite.go:32: setting up integration test suite
    network.go:173: acquiring test network lock
    network.go:178: created temporary directory: /tmp/TestIntegrationTestSuite1190615240/001/chain-IKdhtE2426805842
    network.go:187: preparing test network...
    network.go:381: starting test network...
[...]
    network.go:386: started test network
=== RUN   TestIntegrationTestSuite/TestMyNewTestCase
=== RUN   TestIntegrationTestSuite/TestCLPsExists
=== RUN   TestIntegrationTestSuite/TestRowanBalanceExists
=== CONT  TestIntegrationTestSuite
    suite.go:71: tearing down integration test suite
    network.go:473: cleaning up test network...
    network.go:496: finished cleaning up test network
    network.go:470: released test network lock
--- PASS: TestIntegrationTestSuite (16.46s)
    --- PASS: TestIntegrationTestSuite/TestCLPsExists (0.00s)
    --- PASS: TestIntegrationTestSuite/TestRowanBalanceExists (0.00s)
    --- PASS: TestMyNewTestCase (0.00s)
PASS
ok      github.com/Sifchain/sifnode/x/pmtp/client/testutil      (cached)
```
