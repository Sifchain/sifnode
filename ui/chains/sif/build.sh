#!/bin/bash

killall sifnoded sifnodecli

rm $(which sifnoded) 2> /dev/null || echo sifnoded not install yet ...
rm $(which sifnodecli) 2> /dev/null || echo sifnodecli not install yet ...

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

cd ../../../ && make install 