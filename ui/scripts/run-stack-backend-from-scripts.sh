#!/bin/bash


# This will run our stack from genesis
# All setup scripts will be run and once it is ready our app will be compiled and built
# You probably only need to use this script to set a new snapshot archive

# Run a tmux session with all the background services running requires tmux
if ! command -v tmux &> /dev/null
then
    echo "You need tmux to run this script. Install instructions: https://github.com/tmux/tmux/wiki"
    exit
fi

killall sifnoded sifnodecli ganache-cli

#!/bin/bash

yarn concurrently -k -r -s first \
  "yarn chain:eth" \
  "yarn chain:sif" \
  "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 && yarn chain:migrate && yarn chain:peggy"
