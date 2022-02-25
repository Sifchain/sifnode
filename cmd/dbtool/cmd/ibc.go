package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/spf13/cobra"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func NewIBCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc",
		Short: "IBC queries",
	}
	cmd.AddCommand(
		pendingTransfersCmd,
		connectionCmd,
		getTransfersCmd,
	)
	return cmd
}

var pendingTransfersCmd = &cobra.Command{
	Use:   "pending channel-id",
	Short: "Get pending transfers",
	Long: `Get pending transfers
	
Return the list of fungible token packets that haven't been acknowledged or
timed out`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getPendingTxs(args[0])
	},
}

var connectionCmd = &cobra.Command{
	Use:   "connection connection-id",
	Short: "Get connection info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getConnection(args[0])
	},
}

var getTransfersCmd = &cobra.Command{
	Use:   "get-transfers input_file channel-id node-url",
	Short: "get send_packet data by packet sequence",
	Long: `
Read a list of packet sequences and fetch the corresponding send_packet from the
specified node and channel-id.

ATTENTION: This command does not read from the database

ex: dbtool ibc get-transfers ~/stuck_packets.txt channel-18 http://rpc.sifchain.finance:80
`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		getTransfers(args[0], args[1], args[2])
	},
}

func getPendingTxs(channel string) {
	sifApp, err := utils.NewSifApp(datadir)
	if err != nil {
		panic(err)
	}

	lastBlockHeight := sifApp.LastBlockHeight()

	ctx := sifApp.NewContext(
		true,
		tmproto.Header{Height: lastBlockHeight},
	)

	commitments := sifApp.IBCKeeper.ChannelKeeper.GetAllPacketCommitmentsAtChannel(
		ctx,
		"transfer",
		channel,
	)

	for _, commitment := range commitments {
		query := fmt.Sprintf("send_packet.packet_src_channel='%s' AND send_packet.packet_sequence=%d", channel, commitment.Sequence)
		utils.Print(fmt.Sprintf("query: %s\n", query))

		txs, err := utils.DoTxSearch(
			query,
			true,
			1,
			100,
		)
		if err != nil {
			panic(err)
		}

		filter := func(eventType string) bool {
			return eventType == "send_packet"
		}
		filteredEvents := utils.FilterEvents(txs, filter)

		if len(filteredEvents) == 0 {
			fmt.Printf("%d PRUNED\n", commitment.Sequence)
			continue
		}
		if len(filteredEvents) > 1 {
			panic(fmt.Errorf("Multiple events (%d) for %s", len(filteredEvents), query))
		}

		ev := filteredEvents[0]
		var fungibleTokenPacket FungibleTokenPacketData
		err = json.Unmarshal(
			[]byte(ev.GetAttribute("packet_data")),
			&fungibleTokenPacket,
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s tx_hash:%s, amount:%s, denom: %s, receiver: %s, sender: %s\n",
			ev.GetAttribute("packet_sequence"),
			ev.TxHash,
			fungibleTokenPacket.Amount,
			fungibleTokenPacket.Denom,
			fungibleTokenPacket.Receiver,
			fungibleTokenPacket.Sender,
		)
	}
}

func getConnection(connectionID string) {
	sifApp, err := utils.NewSifApp(datadir)
	if err != nil {
		panic(err)
	}

	lastBlockHeight := sifApp.LastBlockHeight()

	ctx := sifApp.NewContext(
		true,
		tmproto.Header{Height: lastBlockHeight},
	)

	connection, _ := sifApp.IBCKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)

	fmt.Println(connection.String())

	resp, err := sifApp.IBCKeeper.ChannelKeeper.ConnectionChannels(
		sdk.WrapSDKContext(ctx),
		&types.QueryConnectionChannelsRequest{
			Connection: connectionID,
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", resp)
}

func getTransfers(packetFile string, channelID string, node string) {
	file, err := os.Open(packetFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	c, err := rpchttp.New(node, "/websocket")
	if err != nil {
		panic(err)
	}

	transfers := []*Transfer{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		seq, _ := strconv.ParseUint(scanner.Text(), 10, 64)
		ev, err := getSendEvent(c, uint64(seq), channelID)
		if err != nil {
			panic(err)
		}
		transfers = append(transfers, ev)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// print results in csv format
	fmt.Println("packet_sequence, tx_hash, amount, denom, receiver, sender")
	for _, t := range transfers {
		fmt.Printf("%s, %s, %s, %s, %s, %s\n",
			t.Sequence,
			t.TxHash,
			t.PacketData.Amount,
			t.PacketData.Denom,
			t.PacketData.Receiver,
			t.PacketData.Sender,
		)
	}
}

func getSendEvent(client *rpchttp.HTTP, packetSequence uint64, channelID string) (*Transfer, error) {
	query := fmt.Sprintf("send_packet.packet_sequence=%d AND send_packet.packet_src_channel='%s'", packetSequence, channelID)

	page := 1
	perPage := 100
	res, err := client.TxSearch(
		context.Background(),
		query,
		false,
		&page,
		&perPage,
		"asc",
	)
	if err != nil {
		return nil, err
	}

	filter := func(eventType string) bool {
		return eventType == "send_packet"
	}
	filteredEvents := utils.FilterEvents(res.Txs, filter)

	if len(filteredEvents) == 0 {
		return nil, fmt.Errorf("Pruned send_packet (sequence %d)", packetSequence)
	}
	if len(filteredEvents) > 1 {
		return nil, fmt.Errorf("Multiple events (%d) for %s", len(filteredEvents), query)
	}

	ev := filteredEvents[0]
	var fungibleTokenPacket FungibleTokenPacketData
	err = json.Unmarshal(
		[]byte(ev.GetAttribute("packet_data")),
		&fungibleTokenPacket,
	)
	if err != nil {
		return nil, err
	}

	t := &Transfer{
		Sequence:   ev.GetAttribute("packet_sequence"),
		TxHash:     ev.TxHash,
		PacketData: fungibleTokenPacket,
	}

	return t, nil
}

type Transfer struct {
	Sequence   string
	TxHash     string
	PacketData FungibleTokenPacketData
}

type FungibleTokenPacketData struct {
	// the token denomination to be transferred
	Denom string `json:"denom,omitempty"`
	// the token amount to be transferred
	Amount string `json:"amount,omitempty"`
	// the sender address
	Sender string `json:"sender,omitempty"`
	// the recipient address on the destination chain
	Receiver string `json:"receiver,omitempty"`
}
