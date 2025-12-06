PROJECT  := toll
SERVICES := api

.DEFAULT_GOAL := help
.PHONY: help test/static test/static/docker lint bump/%

help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

build: $(addprefix build/,$(SERVICES)) ## Build all builder containers
build/%: ## Build a builder container (e.g. build/api)
	docker build \
		--target builder \
		--tag $(PROJECT)/$*:latest \
		--file $(CURDIR)/$*/Dockerfile \
		../..
