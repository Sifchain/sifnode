#!/usr/bin/env bash
rm -rf all.json pending.json completed.json
sifnodecli q dispensation records-by-name-all ar1 >> all.json
sifnodecli q dispensation records-by-name-pending ar1 >> pending.json
sifnodecli q dispensation records-by-name-completed ar1 >> completed.json