.DEFAULT_GOAL := help

# AutoDoc
# -------------------------------------------------------------------------
.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: gotidy
gotidy: ## Run golangci-lint, goimports and gofmt
	golangci-lint run ./... && goimports -w  . && gofmt -s -w -e -d .

.PHONY: gotest
gotest: ## Run integration and unit tests
	go test ./... -cover -coverpkg=./... --tags=unit,integration

.PHONY: generate
generate: ## Generate and install a new version of kli
	./scripts/generate_release.sh 1.0.0 && sudo mv dist/kli /usr/bin/kli

