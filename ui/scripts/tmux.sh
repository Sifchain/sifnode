#!/bin/bash

# Run a tmux session with all the background services running requires tmux

tmux \
  new-session 'yarn app:serve' \; \
  split-window 'yarn chain:eth' \; \
  split-window 'yarn core:watch' \; \
  select-layout even-horizontal \; \
  select-pane -L \; \
  select-pane -L \; \
  split-window  'echo "wait on websocket..." && sleep 5 && wscat -c "ws://localhost:26657/websocket" -x "{ \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"tm.event='"'"'Tx'"'"'\"], \"id\": 1 }" -w 4000' \; \
  select-pane -R \; \
  split-window 'yarn chain:sif' \; \
  select-pane -R \; \
  split-window 'yarn chain:migrate && yarn chain:peggy' \; \
  
  
