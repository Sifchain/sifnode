package ibctransfer_test_transfer

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

	// testing chains used for convenience and readability
	chainA *TestChain
	chainB *TestChain
	chainC *TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(GetChainID(2))
}

func (suite *TransferTestSuite) TestSingleTransfer() {
	timeoutHeight := clienttypes.NewHeight(0, 110)
	testDenom := "cphoton"
	testAmount := int64(10123456789)
	testCoin := sdk.NewCoin(testDenom, sdk.NewInt(testAmount))

	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)

	// send from chainA to chainB
	msg := transfertypes.NewMsgTransfer(channelA.PortID, channelA.ID, testCoin, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err)

	// relay send
	fungibleTokenPacket := transfertypes.NewFungibleTokenPacketData(testCoin.Denom, testCoin.Amount.Uint64(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check that voucher exists on chain B
	voucherDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), testDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	coinSentFromAToB := transfertypes.GetTransferCoin(channelB.PortID, channelB.ID, testDenom, testAmount)
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(99989876543211), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Equal(coinSentFromAToB, balance)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
}

func (suite *TransferTestSuite) TestMultiTransfer() {
	timeoutHeight := clienttypes.NewHeight(0, 110)
	testDenom := "cusdt"
	testAmount := int64(10123456789)
	testCoin := sdk.NewCoin(testDenom, sdk.NewInt(testAmount))

	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)

	// send from chainA to chainB
	msg := transfertypes.NewMsgTransfer(channelA.PortID, channelA.ID, testCoin, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err)

	// relay send
	fungibleTokenPacket := transfertypes.NewFungibleTokenPacketData(testCoin.Denom, testCoin.Amount.Uint64(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check that voucher exists on chain B
	voucherDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), testDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	coinSentFromAToB := transfertypes.GetTransferCoin(channelB.PortID, channelB.ID, testDenom, testAmount)
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(99989876543211), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Equal(coinSentFromAToB, balance)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)

	// setup between chainB to chainC
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// send from chainB to chainC
	msg = transfertypes.NewMsgTransfer(channelOnBForC.PortID, channelOnBForC.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainC, clientOnCForB, msg)
	suite.Require().NoError(err)

	// RELAY SEND: Fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPath := transfertypes.GetPrefixedDenom(channelOnCForB.PortID, channelOnCForB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.Uint64(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	coinSentFromBToC := sdk.NewInt64Coin(transfertypes.ParseDenomTrace(fullDenomPath).IBCDenom(), testAmount)

	// check that the balance on chain C is correct
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Equal(coinSentFromBToC, balance)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// check that balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// check that balance on chain A is empty
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(99989876543211), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
}

func (suite *TransferTestSuite) TestSingleTransferReverse() {
	timeoutHeight := clienttypes.NewHeight(0, 110)
	testDenom := "cphoton"
	testAmount := int64(10123456789)
	testCoin := sdk.NewCoin(testDenom, sdk.NewInt(testAmount))

	// setup between chain A and chain B
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)

	// send from chain A to chain B
	msg := transfertypes.NewMsgTransfer(channelA.PortID, channelA.ID, testCoin, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err)

	// relay send
	fungibleTokenPacket := transfertypes.NewFungibleTokenPacketData(testCoin.Denom, testCoin.Amount.Uint64(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check that voucher exists on chain B
	voucherDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), testDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	coinSentFromAToB := transfertypes.GetTransferCoin(channelB.PortID, channelB.ID, testDenom, testAmount)
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(99989876543211), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Equal(coinSentFromAToB, balance)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)

	// send from chain B back to chain A
	msg = transfertypes.NewMsgTransfer(channelB.PortID, channelB.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainA.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, clientA, msg)
	suite.Require().NoError(err)

	// RELAY SEND: Fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPathBack := transfertypes.GetPrefixedDenom(channelB.PortID, channelB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.Uint64(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelB.PortID, channelB.ID, channelA.PortID, channelA.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainA, clientB, clientA, packet, ack.GetBytes())
	suite.Require().NoError(err)

	coinSentFromBToA := sdk.NewInt64Coin(transfertypes.ParseDenomTrace(fullDenomPathBack).IBCDenom(), testAmount)

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToA.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToA.Denom)
	suite.Require().Zero(balance.Amount.Int64())
}

func (suite *TransferTestSuite) TestMultiTransferReverse() {
	timeoutHeight := clienttypes.NewHeight(0, 110)
	testDenom := "ceth"
	testAmount := int64(59876543210405)
	testCoin := sdk.NewCoin(testDenom, sdk.NewInt(testAmount))

	// setup between chainA and chain B
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)

	// send from chainA to chain B
	msg := transfertypes.NewMsgTransfer(channelA.PortID, channelA.ID, testCoin, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	fungibleTokenPacket := transfertypes.NewFungibleTokenPacketData(testCoin.Denom, testCoin.Amount.Uint64(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	voucherDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), testDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	coinSentFromAToB := transfertypes.GetTransferCoin(channelB.PortID, channelB.ID, testDenom, testAmount)
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(40123456789595), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Equal(coinSentFromAToB, balance)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Equal(sdk.NewInt(testAmount), balance.Amount)

	// check that the balance on chain C is correct
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// setup between chainB to chain C
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// send from chainB to chain C
	msg = transfertypes.NewMsgTransfer(channelOnBForC.PortID, channelOnBForC.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainC, clientOnCForB, msg)
	suite.Require().NoError(err)

	// RELAY SEND: Fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPath := transfertypes.GetPrefixedDenom(channelOnCForB.PortID, channelOnCForB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.Uint64(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	coinSentFromBToC := sdk.NewInt64Coin(transfertypes.ParseDenomTrace(fullDenomPath).IBCDenom(), testAmount)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(40123456789595), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain C is correct
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Equal(coinSentFromBToC, balance)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// send from chain C back to chain B
	msg = transfertypes.NewMsgTransfer(channelOnCForB.PortID, channelOnCForB.ID, coinSentFromBToC, suite.chainC.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err = suite.coordinator.SendMsg(suite.chainC, suite.chainB, clientOnBForC, msg)
	suite.Require().NoError(err)

	// RELAY SEND: Fungible token is prefixed with the full trace in order to verify the packet commitment
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(fullDenomPath, coinSentFromBToC.Amount.Uint64(), suite.chainC.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnCForB.PortID, channelOnCForB.ID, channelOnBForC.PortID, channelOnBForC.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainC, suite.chainB, clientOnCForB, clientOnBForC, packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(40123456789595), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B returned back to the original state
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Equal(coinSentFromAToB, balance)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Equal(sdk.NewInt(testAmount), balance.Amount)

	// check that module account escrow address is empty
	escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), escrowAddress, testDenom)
	suite.Require().Equal(sdk.NewCoin(testDenom, sdk.ZeroInt()), balance)

	// check that balance on chain C is correct
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// send from chain B back to chain A
	msg = transfertypes.NewMsgTransfer(channelB.PortID, channelB.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainA.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, clientA, msg)
	suite.Require().NoError(err)

	// RELAY SEND: Fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPathBack := transfertypes.GetPrefixedDenom(channelB.PortID, channelB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = transfertypes.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.Uint64(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelB.PortID, channelB.ID, channelA.PortID, channelA.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainA, clientB, clientA, packet, ack.GetBytes())
	suite.Require().NoError(err)

	coinSentFromBToA := sdk.NewInt64Coin(transfertypes.ParseDenomTrace(fullDenomPathBack).IBCDenom(), testAmount)

	// check that balance on chain C is correct
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromAToB.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that the balance on chain B is correct
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToA.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())

	// check that balance on chain A is correct
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), testDenom)
	suite.Require().Equal(sdk.NewInt(100000000000000), balance.Amount)
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), coinSentFromBToA.Denom)
	suite.Require().Zero(balance.Amount.Int64())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
