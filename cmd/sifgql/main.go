package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
)

var (
	flagAppRPC        = flag.String("app-rpc", "localhost:1317", "Application's RPC endpoint.")
	flagTendermintRPC = flag.String("tendermint-rpc", "localhost:26657", "Tendermint's RPC endpoint.")
	flagBlockchain    = flag.String("blockchain", "sifchain", "Application's name (e.g. Cosmos Hub)")
	flagNetworkID     = flag.String("network", "localnet", "Network's identifier (e.g. cosmos-hub-3, testnet-1, etc)")
	// flagOfflineMode   = flag.Bool("offline", false, "Flag that forces the rosetta service to run in offline mode, some endpoints won't work.")
	// flagAddrPrefix    = flag.String("prefix", "sif", "Bech32 prefix of address (e.g. cosmos, iaa, xrn:)")
	flagPort = flag.Uint("port", 8081, "The port where the service is exposed.")
)

const defaultPort = "8081"

// func main() {
// 	flag.Parse()
//
// 	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))
//
// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)
//
// 	log.Printf("connect to http://localhost:%v/ for GraphQL playground", *flagPort)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(int(*flagPort))), nil))
// }

// func main() {
// 	flag.Parse()
//
// 	url := "ws://127.0.0.1:26657/websocket"
//
// 	ws, err := websocket.Dial(url,"","http://localhost/")
// 	if err != nil {
// 		log.Fatal("Err0: ", err)
// 	}
//
// 	defer ws.Close()
//
// 	c := jsonrpc.NewClient(ws)
//
// 	var reply interface{}
//
// 	args := struct {
// 		id string `json:"id"`
// 		params struct {
// 			query string `json:"query"`
// 		} `json:"params"`
// 	}{
// 		id: "0",
// 		params: struct {
//     query string `json:"query"`
// }{
// 			query: "tm.event='NewBlock'",
// 		},
// 	}
//
// 	spew.Dump(args)
//
// 	err = c.Call("subscribe", args, &reply)
// 	if err != nil {
// 		log.Fatal("Err: ", err)
// 	}
//
// 	spew.Dump(reply)
// }

func main() {
	client, err := http.New("tcp://127.0.0.1:26657", "/websocket")
	if err != nil {
		log.Fatal("Err0: ", err)
	}

	err = client.Start()
	if err != nil {
		log.Fatal("err1: ", err)
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := "tm.event='NewBlock'"
	reply, err := client.Subscribe(ctx, "test-client", query)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for e := range reply {
			fmt.Println("got ", e.Data.(types.EventDataNewBlock))
			spew.Dump(e)
		}
	}()

	// <- reply
	time.Sleep(time.Second * 60)
}
