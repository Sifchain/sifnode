#!/bin/bash
pid=$(ps aux | grep "ganache-cli" | grep -v grep | awk '{print $2}')

if [[ ! -z "$pid" ]]; then 
  kill -9 $pid
fi
