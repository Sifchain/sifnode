package types_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestValAddress = "cosmosvaloper1s3uh0vhxr6thj2w0ea2h53j8ra7cxm63cu29mf"
)

var (
	ExpectedKey        = []byte{0x6, 0x8, 0x1, 0x12, 0x14, 0x84, 0x79, 0x77, 0xb2, 0xe6, 0x1e, 0x97, 0x79, 0x29, 0xcf, 0xcf, 0x55, 0x7a, 0x46, 0x47, 0x1f, 0x7d, 0x83, 0x6f, 0x51}
	KeyWithWrongPrefix = []byte{0x1, 0x0}
	InvalidKey         = []byte{0x6}
)

func Test_GetWitnessLockBurnSequencePrefix(t *testing.T) {
	app := sifapp.Setup(false)

	valAddress, _ := sdk.ValAddressFromBech32(TestValAddress)

	cdc := codec.BinaryCodec(app.AppCodec())

	key := types.LockBurnSequenceKey{
		NetworkDescriptor: 1,
		ValidatorAddress:  valAddress,
	}
	value := key.GetWitnessLockBurnSequencePrefix(cdc)
	assert.Equal(t, value, ExpectedKey)
}

func Test_GetWitnessLockBurnSequenceKeyFromRawKey(t *testing.T) {
	app := sifapp.Setup(false)

	valAddress, _ := sdk.ValAddressFromBech32(TestValAddress)

	cdc := codec.BinaryCodec(app.AppCodec())

	val, err := types.GetWitnessLockBurnSequenceKeyFromRawKey(cdc, ExpectedKey)
	assert.NoError(t, err)
	assert.Equal(t, val.ValidatorAddress, valAddress.Bytes())

	val, err = types.GetWitnessLockBurnSequenceKeyFromRawKey(cdc, KeyWithWrongPrefix)
	assert.Error(t, err, "LockBurnSequenceKey prefix is invalid")

	// even the invalid data marshal without the error
	val, err = types.GetWitnessLockBurnSequenceKeyFromRawKey(cdc, InvalidKey)
	assert.Error(t, err)
}
