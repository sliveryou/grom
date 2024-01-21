.PHONY: build install pkg clean fmt proxy

build:
	@sh scripts/build.sh

install:
	@sh scripts/install.sh

pkg:
	@sh scripts/package.sh

clean:
	@sh scripts/clean.sh

fmt:
	@find . -name '*.go' -not -path "./vendor/*" | xargs gofumpt -w -extra
	@find . -name '*.go' -not -path "./vendor/*" | xargs -n 1 -t goimports-reviser -rm-unused -set-alias -company-prefixes "github.com/sliveryou" -project-name "github.com/sliveryou/grom"
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shfmt -w -s -i 2 -ci -bn -sr

proxy:
	@go env -w GO111MODULE="on"
	@go env -w GOPROXY="https://goproxy.cn,direct"
