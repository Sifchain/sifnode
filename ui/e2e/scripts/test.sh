#!/bin/bash

set -e

if ! command -v xvfb-run &> /dev/null
then
  yarn test
else 
  xvfb-run --auto-servernum -- yarn test
fi