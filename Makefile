build:
	@sh scripts/build.sh

install:
	@sh scripts/install.sh

pkg:
	@sh scripts/package.sh

clean:
	@sh scripts/clean.sh

fmt:
	@find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -s -w
	@find . -name '*.go' -not -path "./vendor/*" | xargs goimports -w
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shfmt -w -s -i 4 -ci -bn

proxy:
	@go env -w GO111MODULE="on"
	@go env -w GOPROXY="https://goproxy.cn,direct"

dep:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
	@go get golang.org/x/tools/cmd/goimports
	@go get mvdan.cc/sh/v3/cmd/shfmt
	@go get mvdan.cc/sh/v3/cmd/gosh
	@git checkout -- go.mod go.sum
