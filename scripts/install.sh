#!/bin/sh
file_path=$(
  cd "$(dirname "$0")" || exit
  pwd
)/..

main_path="github.com/sliveryou/grom/cmd"
go_version=$(go version | awk '{ print $3 }')
build_time=$(date "+%Y-%m-%d %H:%M:%S")
git_commit=$(git rev-parse --short=10 HEAD)
flags="-X '${main_path}.goVersion=${go_version}' -X '${main_path}.buildTime=${build_time}' -X '${main_path}.gitCommit=${git_commit}'"

cd "${file_path}" || exit
go install -ldflags "${flags}"
