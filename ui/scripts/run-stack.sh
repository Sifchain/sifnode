#!/bin/bash

# This will run our stack from the stored snapshot

# Run a tmux session with all the background services running requires tmux
if ! command -v tmux &> /dev/null
then
    echo "You need tmux to run this script. Install instructions: https://github.com/tmux/tmux/wiki"
    exit
fi

killall sifnoded sifnodecli ebrelayer ganache-clir

tmux \
  new-session 'yarn stack:backend || echo "process finished" && sleep 1000' \; \
  split-window 'yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 && yarn app:serve || echo "process finished" && sleep 1000' \; \
  split-window 'yarn core:watch || echo "process finished" && sleep 1000' \; \
  select-layout even-vertical

## CHEAT SHEET

# ctrl+b (left|right|up|down) - select window
# ctrl+b [ - scroll mode 
# ctrl+c - exit 