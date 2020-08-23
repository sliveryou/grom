#!/bin/sh
file_path=$(
    cd $(dirname $0)
    pwd
)/..

rm -r ${file_path}/grom*
