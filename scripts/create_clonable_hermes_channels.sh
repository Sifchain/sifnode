#!/bin/zsh

hermes create channel localnet-1 localnet-2 --port-a transfer --port-b transfer -o unordered
hermes create channel localnet-1 localnet-3 --port-a transfer --port-b transfer -o unordered
