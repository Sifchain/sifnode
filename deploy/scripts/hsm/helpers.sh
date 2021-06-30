#!/usr/bin/env bash
set +x

function header (){
    echo "======================"
    echo "${1}"
    echo "======================"
}

function footer (){
    echo "======================"
    echo ""
}

function log(){
    echo " "
    echo $(date) - $1
    echo " "
}

function get_cluster_id(){
    echo " "
    echo $(date) - $1
    echo " "
}