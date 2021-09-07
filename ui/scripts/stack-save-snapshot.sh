#!/bin/bash

# Run the stack until migration is complete and then take a snapshot

# reset our migrate complete flag see chains/post_migrate.sh
rm node_modules/.migrate-complete

yarn concurrently -r \
  "./scripts/stack-launch.sh" \
  "yarn wait-on http-get://localhost:1317/node_info node_modules/.migrate-complete && sleep 30 && ./scripts/_snapshot.sh"