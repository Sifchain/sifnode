#!/usr/bin/env bash
rm -rf all.json pending.json completed.json
sifnoded q dispensation records-by-name ar1 All>> all.json
sifnoded q dispensation records-by-name ar1 Pending >> pending.json
sifnoded q dispensation records-by-name ar1 Completed>> completed.json