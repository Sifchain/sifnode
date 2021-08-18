package test

import (
	"testing"

	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/stretchr/testify/suite"
)

type TransferTestSuite struct {
	suite.Suite
	coordinator *Coordinator
	t           *testing.T

	// testing chains used for convenience and readability
	chainA *TestChain
	chainB *TestChain
	chainC *TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	_, suite.chainA = CreateTestChain(suite.t, "1", false)
	_, suite.chainB = CreateTestChain(suite.t, "2", false)
	_, suite.chainC = CreateTestChain(suite.t, "3", false)

	chains := make(map[string]*TestChain)
	chains[suite.chainA.ChainID] = suite.chainA
	chains[suite.chainB.ChainID] = suite.chainB
	chains[suite.chainC.ChainID] = suite.chainC

	suite.coordinator = NewCoordinator(suite.t, chains)
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)
	timeoutHeight := clienttypes.NewHeight(0, 110)
	//testDenom := "cusdt"
	testDenom := sdk.DefaultBondDenom
	coinToSendToB := sdk.NewCoin(testDenom, sdk.NewInt(100))
	// initCoins := sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(555)))

	// err := suite.chainA.App.BankKeeper.AddCoins(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), initCoins)
	// if err != nil {
	// 	panic(err)
	// }

	// err = suite.chainB.App.BankKeeper.AddCoins(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), initCoins)
	// if err != nil {
	// 	panic(err)
	// }

	// send from chainA to chainB
	msg := transfertypes.NewMsgTransfer(channelA.PortID, channelA.ID, coinToSendToB, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)

	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	fungibleTokenPacket := transfertypes.NewFungibleTokenPacketData(coinToSendToB.Denom, coinToSendToB.Amount.Uint64(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	voucherDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), testDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())

	coinSentFromAToB := transfertypes.GetTransferCoin(channelB.PortID, channelB.ID, testDenom, 100)
	suite.Require().Equal(coinSentFromAToB, balance)

	// setup between chainB to chainC
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// send from chainB to chainC
	msg = transfertypes.NewMsgTransfer(channelOnBForC.PortID, channelOnBForC.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0)

	err = suite.coordinator.SendMsg(suite.chainB, suite.chainC, clientOnCForB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPath := transfertypes.GetPrefixedDenom(channelOnCForB.PortID, channelOnCForB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.Uint64(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	coinSentFromBToC := sdk.NewInt64Coin(transfertypes.ParseDenomTrace(fullDenomPath).IBCDenom(), 100)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)

	// check that the balance is updated on chainC
	suite.Require().Equal(coinSentFromBToC, balance)

	// check that balance on chain B is empty
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// send from chainC back to chainB
	msg = transfertypes.NewMsgTransfer(channelOnCForB.PortID, channelOnCForB.ID, coinSentFromBToC, suite.chainC.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)

	err = suite.coordinator.SendMsg(suite.chainC, suite.chainB, clientOnBForC, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(fullDenomPath, coinSentFromBToC.Amount.Uint64(), suite.chainC.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnCForB.PortID, channelOnCForB.ID, channelOnBForC.PortID, channelOnBForC.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainC, suite.chainB, clientOnCForB, clientOnBForC, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)

	// check that the balance on chainA returned back to the original state
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that module account escrow address is empty
	escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), escrowAddress, testDenom)
	suite.Require().Equal(sdk.NewCoin(testDenom, sdk.ZeroInt()), balance)

	// check that balance on chain B is empty
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())
}

func TestTransferTestSuite(t *testing.T) {
	testingSuite := &TransferTestSuite{
		t: t,
	}

	suite.Run(t, testingSuite)
}
