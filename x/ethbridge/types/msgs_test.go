package types

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/CosmosBridge.sol"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// TODO: Update this value with the expected hash after fixing TestComputeTest
var expectedProphecyID = []byte{0x11, 0x96, 0x5a, 0xd2, 0x8e, 0xe1, 0xf5, 0xdd, 0x50, 0xb5, 0xef, 0xcb, 0x6c, 0xb2, 0xa4, 0x7d, 0xf5, 0x55, 0x8b, 0x5b, 0xc5, 0x9b, 0x48, 0x51, 0x76, 0xf7, 0x67, 0x5c, 0xf, 0x85, 0xf3, 0xce}

// Test that verifies compute prophecy works as expected AND that every field of claim data is correctly hashed in the prophecy ID
func TestComputeProphecyID(t *testing.T) {
	// Create a claim data object to verify all fields in claim data are populated
	claimData := cosmosbridge.CosmosBridgeClaimData{
		CosmosSender:         []byte("cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"),
		CosmosSenderSequence: big.NewInt(0),
		EthereumReceiver:     common.HexToAddress("0xa98cea040E91e28D71b883b88d6c6445b486124D"),
		TokenAddress:         common.HexToAddress("0xC62C770B3223E7ABeD54B4026ad972C84e9a424b"),
		Amount:               big.NewInt(1025),
		TokenName:            "Test Token",
		TokenSymbol:          "TT",
		TokenDecimals:        18,
		NetworkDescriptor:    1,
		BridgeToken:          true,
		Nonce:                big.NewInt(100),
		CosmosDenom:          "sifBridge0123456789",
	}

	// Specify the abi types
	bytesTy, _ := abi.NewType("bytes", "bytes", nil)
	boolTy, _ := abi.NewType("bool", "bool", nil)
	uint8Ty, _ := abi.NewType("uint8", "uint8", nil)
	int32Ty, _ := abi.NewType("int32", "int32", nil)
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	addressTy, _ := abi.NewType("address", "address", nil)
	stringTy, _ := abi.NewType("string", "string", nil)

	// Specify the abi packing layout
	arguments := abi.Arguments{
		{
			Type: bytesTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: stringTy,
		},
		{
			Type: stringTy,
		},
		{
			Type: uint8Ty,
		},
		{
			Type: int32Ty,
		},
		{
			Type: boolTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: stringTy,
		},
	}

	// Iterate over every field in claimData, pack it in the abi format and hash the result
	res := reflect.ValueOf(claimData)
	packedSlice := make([]interface{}, res.NumField())
	for i := 0; i < res.NumField(); i++ {
		fmt.Println(res.Field(i))
		packedSlice[i] = res.Field(i).Interface()
	}
	packedData, err := arguments.Pack(packedSlice...)
	assert.NoError(t, err)
	hashBytes := crypto.Keccak256(packedData)

	// Compute the prophecy ID
	prophecy := ComputeProphecyID(
		string(claimData.CosmosSender),
		claimData.CosmosSenderSequence.Uint64(),
		claimData.EthereumReceiver.Hex(),
		claimData.TokenAddress.Hex(),
		sdk.NewIntFromBigInt(claimData.Amount),
		claimData.TokenName,
		claimData.TokenSymbol,
		claimData.TokenDecimals,
		oracletypes.NetworkDescriptor(claimData.NetworkDescriptor),
		claimData.BridgeToken,
		claimData.Nonce.Uint64(),
		claimData.CosmosDenom)

	// Verify the hash of all fields of claim data matches the prophecy ID of compute prophecyID
	assert.Equal(t, hashBytes, prophecy)
	// Verify the hash of prophecy ID matches the expected hash value
	assert.Equal(t, expectedProphecyID, prophecy)
}
