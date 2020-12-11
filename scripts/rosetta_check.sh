#!/usr/bin/env bash

make rosetta &

curl --location --request POST 'http://localhost:8080/network/status' --header 'Content-Type: text/plain' --data-raw '{
    "network_identifier": {
        "blockchain": "sifchain",
        "network": "monkey-bars-local"
    }
}'
