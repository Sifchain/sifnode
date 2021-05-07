#!/usr/bin/env bash

killall sifnoded

rm $(which sifnoded) 2> /dev/null || echo sifnoded not install yet ...

rm -rf ~/.sifnoded

cd ../../../ && make install 