package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"
	"github.com/spf13/cobra"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func NewIBCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc",
		Short: "IBC queries",
	}
	cmd.AddCommand(
		pendingTransfersCmd,
		connectionCmd,
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
			panic(fmt.Errorf("No events for %s", query))
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
		fmt.Printf("tx_hash:%s, packet_sequence:%s, amount:%s, denom: %s, receiver: %s, sender: %s\n",
			ev.TxHash,
			ev.GetAttribute("packet_sequence"),
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
