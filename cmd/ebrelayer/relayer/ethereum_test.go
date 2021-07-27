package relayer

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

// EventProcessed check if the event processed by relayer
func TestEventProcessed(t *testing.T) {
	var bridgeClaims []types.EthereumBridgeClaim
	valAddress, _ := sdk.ValAddressFromBech32("sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy")
	bridgeClaims = append(bridgeClaims, types.EthereumBridgeClaim{
		EthereumSender: common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		CosmosSender:   valAddress,
		Nonce:          sdk.NewInt(int64(1)),
	})

	processedEvent := types.EthereumEvent{
		EthereumChainID:       big.NewInt(0),
		BridgeContractAddress: common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		ID:                    [32]byte{},
		From:                  common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		To:                    []byte("sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy"),
		Token:                 common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		Symbol:                "rewan",
		Value:                 big.NewInt(1),
		Nonce:                 big.NewInt(1),
		ClaimType:             ethbridge.ClaimType_CLAIM_TYPE_LOCK,
	}

	require.Equal(t, true, EventProcessed(bridgeClaims, processedEvent))

	notProcessedEvent := types.EthereumEvent{
		EthereumChainID:       big.NewInt(0),
		BridgeContractAddress: common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		ID:                    [32]byte{},
		From:                  common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		To:                    []byte("sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy"),
		Token:                 common.HexToAddress("0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"),
		Symbol:                "rewan",
		Value:                 big.NewInt(1),
		Nonce:                 big.NewInt(10),
		ClaimType:             ethbridge.ClaimType_CLAIM_TYPE_LOCK,
	}

	require.Equal(t, false, EventProcessed(bridgeClaims, notProcessedEvent))
}

// TestNewKeybase test if we can get keybase from moniker, mnemonic and password
func TestNewKeybase(t *testing.T) {
	validatorMoniker := "akasha"
	mnemonic := "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
	password := ""

	base, info, err := NewKeybase(validatorMoniker, mnemonic, password)
	require.NotEqual(t, base, nil)
	require.NotEqual(t, info, nil)
	require.Equal(t, err, nil)
}
