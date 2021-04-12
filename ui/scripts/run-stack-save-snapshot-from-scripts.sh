#!/bin/bash

# Run the stack until migration is complete and then take a snapshot

# reset our migrate complete flag
rm node_modules/.migrate-complete

yarn concurrently -r \
  "yarn stack:backend-from-scripts" \
  "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 node_modules/.migrate-complete && sleep 10 && yarn chain:snapshot"