export PROJECTS_ROOT := $(abspath $(shell pwd)/../)

ensure-bin:
	[ -d .bin ] || mkdir .bin

setup: ensure-bin ## Install tools
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.27.0
	mv bin/golangci-lint .bin/golangci-lint && rm -rf bin

install-tools: install-golangci-lint ## install tools

install-golangci-lint: ensure-bin 
	cd .bin && \
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.17.1  && \
	mv bin/golangci-lint golangci-lint && rm -rf bin

lint: ## Run the linters
	golangci-lint run

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