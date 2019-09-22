export PROJECTS_ROOT := $(abspath $(shell pwd)/../)

ensure-bin:
	[ -d .bin ] || mkdir .bin
	
install-tools: install-golangci-lint ## install tools

install-golangci-lint: ensure-bin 
	cd .bin && \
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.17.1  && \
	mv bin/golangci-lint golangci-lint && rm -rf bin

lint: ## run all the linters
	golangci-lint help linters
	golangci-lint run --enable=staticcheck --enable=govet --enable=interfacer --enable=scopelint --enable=misspell --enable=depguard --enable=dupl --enable=goimports --enable=gofmt --enable=gocyclo --enable=nakedret --enable=scopelint --enable=stylecheck --disable=typecheck
	
tests: ## run all the tests
	go version
	go env
	go list ./... | xargs -n1 -I{} sh -c 'go test -race {}'

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
	
# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help