#!/bin/sh
file_path=$(
  cd "$(dirname "$0")" || exit
  pwd
)/..

rm -r "${file_path}"/grom*
