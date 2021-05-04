#!/usr/bin/env bash
rm -rf all.json pending.json completed.json
sifnodecli q dispensation records-by-name ar1 All>> all.json
sifnodecli q dispensation records-by-name ar1 Pending >> pending.json
sifnodecli q dispensation records-by-name ar1 Completed>> completed.json