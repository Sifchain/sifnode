#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure

basedir=$(dirname $0)
pidfile=${basedir}/watcher_pids.txt

if [ -f $pidfile ]
then
  pids=$(cat $pidfile)
  for i in $pids
  do
    kill -9 $i || true
  done
  rm $pidfile
fi
