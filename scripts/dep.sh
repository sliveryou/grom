#!/bin/bash
export GO111MODULE=on

download=""
go_version=$(go version | awk '{ print $3 }')
# go 1.17 之后下载编译成可执行文件要使用 go install
if ! (printf '%s\n%s\n' "go1.17" "${go_version}" | sort -V -C); then
  echo "go get dependencies..."
  download="go get"
else
  echo "go install dependencies..."
  download="go install"
fi

# jq: https://github.com/jqlang/jq/releases - mac: `brew install jq` or `brew upgrade jq`
# protoc: https://github.com/protocolbuffers/protobuf/releases/tag/v25.1 - v4.25.1

${download} google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
${download} google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
${download} github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
${download} github.com/incu6us/goimports-reviser/v3@v3.6.2
${download} mvdan.cc/gofumpt@v0.5.0
${download} mvdan.cc/sh/v3/cmd/shfmt@v3.7.0
${download} mvdan.cc/sh/v3/cmd/gosh@v3.7.0

echo "done"
