package keeper_test

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

const networkID = uint32(1)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)
	networkDescriptor := types.NewNetworkDescriptor(networkID)
	whilelist := types.ValidatorWhiteList{WhiteList: make(map[string]uint32)}
	for _, address := range valAddresses {
		fmt.Printf("address is %s\n", address.String())
		whilelist.GetWhiteList()[address.String()] = 100
	}

	app.OracleKeeper.SetOracleWhiteList(ctx, networkDescriptor, whilelist)

	vList := app.OracleKeeper.GetOracleWhiteList(ctx, networkDescriptor)
	assert.Equal(t, len(vList.GetAllValidators()), 2)
	assert.True(t, app.OracleKeeper.ExistsOracleWhiteList(ctx, networkDescriptor))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)
	networkDescriptor := types.NewNetworkDescriptor(networkID)
	whitelist := make(map[string]uint32)

	for _, address := range valAddresses {
		fmt.Printf("address is %s\n", address.String())
		whitelist[address.String()] = 100
	}

	app.OracleKeeper.SetOracleWhiteList(ctx, networkDescriptor, types.ValidatorWhiteList{WhiteList: whitelist})
	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[0]))
	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[1]))
	addresses = sifapp.CreateRandomAccounts(3)
	valAddresses = sifapp.ConvertAddrsToValAddrs(addresses)
	assert.False(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[2]))
}
