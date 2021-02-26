#!/bin/bash
# must run from the root directory of the sifnode tree

# add 18 zeros to a number
to_wei () { echo "${1}000000000000000000" ; }
to_json () { ruby -ryaml -rjson -e "puts YAML::load(STDIN.read).to_json"; }
fullpath () { ruby -e "puts File.expand_path(\"$1\")"; }
logecho () {
  date=$(date +I\[%Y-%m-%d\|%H:%M:%S.%N\])
  echo $date $*
}
filenamedate () { date +%Y-%m-%d-%H-%M-%S-%N; }
# sets an environment variable and writes it to the file
set_persistant_env_var () {
  if [[ "$4" = 'required' && -z "$2" ]]; then
    echo environment variable $1 cannot be empty
    exit 1
  fi
  export $1=$2
  echo "export $1=\"$2\"" >> $3
}
