#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure

# add 18 zeros to a number
to_wei () { echo "${1}000000000000000000" ; }
to_json () { ruby -ryaml -rjson -e "puts YAML::load(STDIN.read).to_json"; }
