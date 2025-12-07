PACKAGE := $(shell awk '/^module/ {print $$NF}' go.mod)
SERVICE := toll

OGEN_VERSION=v1.14.0

GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_TAG    ?= $(shell git describe --tags --exact-match 2>/dev/null)
BUILD_ID   ?= manual
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -ldflags " \
	-X $(PACKAGE)/internal/version.revision=$(GIT_COMMIT) \
	-X $(PACKAGE)/internal/version.buildID=$(BUILD_ID) \
	-X $(PACKAGE)/internal/version.buildTime=$(BUILD_TIME) \
	-X $(PACKAGE)/internal/version.branch=$(GIT_BRANCH) \
	-X $(PACKAGE)/internal/version.tag=$(GIT_TAG) "

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

build/%: ## Build an executable
	go build \
		$(LDFLAGS) \
		-o ./bin/$(SERVICE)-$* \
		./cmd/$*

run/%: ## Run the executable
	./bin/$(SERVICE)-$*

devbox: build/api ## Build the executable
	./bin/$(SERVICE)-api --conf ./config.yaml -mock

clean: mock/clean ## Remove artefacts
	rm -rf ./build

generate: ogen/clean ogen/generate mock/generate ## Generate all generated code

ogen/clean:
	find api/restapi -type f -not -name "mock_*.go" | xargs rm -f

ogen/generate: ## Generate server code from OpenAPI spec
	go run github.com/ogen-go/ogen/cmd/ogen@${OGEN_VERSION} \
		--config ogen.yaml \
		--target api/restapi \
		-package restapi \
		--clean \
		../spec/openapi.yaml

test: test/static test/unit ## Run entire QA suite

test/static: ## Run static code analysis tests
	golangci-lint run --verbose

test/unit: ## Unit test the code
	go test -mod=vendor -tags=mock -v -cover -race ./...

mock/generate: ## Generate mocks
	MOCKERY_VERSION=false go run github.com/vektra/mockery/v2@${MOCKERY_VERSION} --config ./.mockery.yaml

mock/clean: ## Remove mocks
	find . \
		-name "mock_*.go" \
		-not -path "./vendor/*" \
	| xargs rm -f
