package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

func TestKeeper_SetNativeToken(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)
	token := "ceth"
	gas := sdk.NewInt(0)
	lockCost := sdk.NewInt(0)
	burnCost := sdk.NewInt(0)

	app.OracleKeeper.SetNativeToken(ctx, networkDescriptor, token, gas, lockCost, burnCost)

	tokenStored, err := app.OracleKeeper.GetNativeToken(ctx, networkDescriptor)
	assert.NoError(t, err)
	assert.Equal(t, token, tokenStored)
}
