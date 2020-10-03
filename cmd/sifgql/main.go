package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/Sifchain/sifnode/api/graph/generated"
	"github.com/Sifchain/sifnode/api/graph/resolvers"
)

var (
	flagPort = flag.Uint("port", 8081, "The port where the service is exposed.")
)

func main() {
	flag.Parse()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%v/ for GraphQL playground", *flagPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(int(*flagPort))), nil))
}
