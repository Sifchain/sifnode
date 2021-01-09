#!/bin/bash

# Run a tmux session with all the background services running requires tmux
if ! command -v tmux &> /dev/null
then
    echo "You need tmux to run this script. Install instructions: https://github.com/tmux/tmux/wiki"
    exit
fi

killall sifnoded sifnodecli ganache-cli



tmux \
  new-session 'yarn app:serve || echo "process finished" && sleep 1000' \; \
  split-window 'yarn chain:eth || echo "process finished" && sleep 1000' \; \
  split-window 'yarn core:watch || echo "process finished" && sleep 1000' \; \
  select-layout even-horizontal \; \
  select-pane -L \; \
  select-pane -L \; \
  split-window  'yarn wait-on http-get://localhost:1317/node_info && wscat -c "ws://localhost:26657/websocket" -x "{ \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"tm.event='"'"'Tx'"'"'\"], \"id\": 1 }" -w 4000 || echo "process finished" && sleep 1000' \; \
  select-pane -R \; \
  split-window 'yarn chain:sif || echo "process finished" && sleep 1000' \; \
  select-pane -R \; \
  split-window 'yarn wait-on http-get://localhost:1317/node_info && yarn chain:migrate && yarn chain:peggy || echo "process finished" && sleep 1000' \; \


## CHEAT SHEET

# ctrl+b (left|right|up|down) - select window
# ctrl+b [ - scroll mode 
# ctrl+c - exit 