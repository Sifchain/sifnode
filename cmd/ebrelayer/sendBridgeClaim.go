package main

import (
	"log"
	"strconv"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"

	ebrelayertype "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// RunSendBridgeClaimCmd executes RelayToCosmos according to hardcoded claim
func RunSendBridgeClaimCmd(cmd *cobra.Command, args []string) error {
	cliContext, err := client.GetClientTxContext(cmd)

	if err != nil {
		return err
	}

	if cliContext.From == "" {
		log.Println("Received empty clientContext.From, needed for validating cosmos transaction. Check if --from flag is set")
		return errors.New("Missing from flag ")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())

	sugaredLogger := logger.Sugar()

	validatorName := args[0]

	valAddr, err := relayer.GetValAddressFromKeyring(txFactory.Keybase(), validatorName)
	if err != nil {
		return err
	}

	nonce, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}

	amountString := args[2]
	// 36690000000000000000000 for CD6E206A21F9B058D2385ACC87C97066FAE29D51282AACF81C2B035A1B7B3F86
	amount, ok := cosmostypes.NewIntFromString(amountString)
	if !ok {
		log.Fatalln("failed to parse amount")
	}
	claimType := ebrelayertype.ClaimType_CLAIM_TYPE_BURN

	var claims []*ebrelayertype.EthBridgeClaim

	// claim for https://www.mintscan.io/sifchain/txs/CD6E206A21F9B058D2385ACC87C97066FAE29D51282AACF81C2B035A1B7B3F86
	claim := &ebrelayertype.EthBridgeClaim{
		EthereumChainId:       1,
		BridgeContractAddress: "0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8",
		// we can start from 0 for 0xD377cFFCc52C16bF6e9840E77F78F42Ddb946568
		Nonce:                nonce,
		Symbol:               "rowan",
		TokenContractAddress: "0x07baC35846e5eD502aA91AdF6A9e7aA210F2DcbE",
		EthereumSender:       "0xD377cFFCc52C16bF6e9840E77F78F42Ddb946568",
		CosmosReceiver:       "sif1yq4gv7c4gu7gp9pzr4lguaphue4yvz6svtm6e7",
		ValidatorAddress:     valAddr.String(),
		Amount:               amount,
		ClaimType:            claimType,
	}
	claims = append(claims, claim)

	return txs.RelayToCosmos(txFactory, claims, cliContext, sugaredLogger)
}
