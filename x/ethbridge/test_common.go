package ethbridge

import (
	"bytes"
	"encoding/hex"
	keeperLib "github.com/Sifchain/sifnode/x/oracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strconv"
	"testing"

	"github.com/Sifchain/sifnode/x/oracle"
)

func CreateTestHandler(
	t *testing.T, consensusNeeded float64, validatorAmounts []int64,
) (sdk.Context, oracle.Keeper, bank.Keeper, supply.Keeper, auth.AccountKeeper, []sdk.ValAddress, sdk.Handler) {
	ctx, oracleKeeper, bankKeeper, supplyKeeper,
		accountKeeper, validatorAddresses := oracle.CreateTestKeepers(t, consensusNeeded, validatorAmounts, ModuleName)
	bridgeAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	supplyKeeper.SetModuleAccount(ctx, bridgeAccount)

	cdc := keeperLib.MakeTestCodec()
	bridgeKeeper := NewKeeper(cdc, supplyKeeper, oracleKeeper)
	handler := NewHandler(accountKeeper, bridgeKeeper, cdc)

	return ctx, oracleKeeper, bankKeeper, supplyKeeper, accountKeeper, validatorAddresses, handler
}

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

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}
