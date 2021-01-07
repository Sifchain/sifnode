#!/bin/bash

# Run a tmux session with all the background services running requires tmux

killall sifnoded sifnodecli ganache-cli

tmux \
  new-session 'yarn app:serve' \; \
  split-window 'yarn chain:eth' \; \
  split-window 'yarn core:watch' \; \
  select-layout even-horizontal \; \
  select-pane -L \; \
  select-pane -L \; \
  split-window  'yarn wait-on http-get://localhost:1317/node_info && wscat -c "ws://localhost:26657/websocket" -x "{ \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"tm.event='"'"'Tx'"'"'\"], \"id\": 1 }" -w 4000' \; \
  select-pane -R \; \
  split-window 'yarn chain:sif' \; \
  select-pane -R \; \
  split-window 'yarn wait-on http-get://localhost:1317/node_info && yarn chain:migrate && yarn chain:peggy' \; \
  
  
## CHEAT SHEET

# ctrl+b (left|right|up|down) - select window
# ctrl+b [ - scroll mode 
# ctrl+c - exit