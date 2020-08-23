#!/bin/sh
file_path=$(
    cd $(dirname $0)
    pwd
)/..

main_path="github.com/sliveryou/grom/cmd"
go_version=$(go version | awk '{ print $3 }')
build_time=$(date "+%Y-%m-%d %H:%M:%S")
git_commit=$(git rev-parse --short=10 HEAD)
flags="-X '${main_path}.goVersion=${go_version}' -X '${main_path}.buildTime=${build_time}' -X '${main_path}.gitCommit=${git_commit}'"

cd ${file_path}
mkdir -p grom-darwin-amd64 grom-linux-amd64 grom-linux-arm64 grom-windows-amd64

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$flags" -o ${file_path}/grom-darwin-amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$flags" -o ${file_path}/grom-linux-amd64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$flags" -o ${file_path}/grom-linux-arm64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$flags" -o ${file_path}/grom-windows-amd64

echo ${file_path}/grom-darwin-amd64 ${file_path}/grom-linux-amd64 ${file_path}/grom-linux-arm64 ${file_path}/grom-windows-amd64 | xargs -n 1 cp ${file_path}/README.md
echo ${file_path}/grom-darwin-amd64 ${file_path}/grom-linux-amd64 ${file_path}/grom-linux-arm64 ${file_path}/grom-windows-amd64 | xargs -n 1 cp ${file_path}/README_zh-CN.md

tar -zcvf ${file_path}/grom-darwin-amd64.tar.gz -C ${file_path} grom-darwin-amd64
tar -zcvf ${file_path}/grom-linux-amd64.tar.gz -C ${file_path} grom-linux-amd64
tar -zcvf ${file_path}/grom-linux-arm64.tar.gz -C ${file_path} grom-linux-arm64
tar -zcvf ${file_path}/grom-windows-amd64.tar.gz -C ${file_path} grom-windows-amd64

rm -r ${file_path}/grom-darwin-amd64
rm -r ${file_path}/grom-linux-amd64
rm -r ${file_path}/grom-linux-arm64
rm -r ${file_path}/grom-windows-amd64
