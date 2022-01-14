#!/bin/zsh

hermes create channel localnet-2 localnet-1 --port-a transfer --port-b transfer -o unordered
hermes create channel localnet-2 localnet-3 --port-a transfer --port-b transfer -o unordered
