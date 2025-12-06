EVENTS_NAME ?= toll
EVENTS_USER ?= toll
EVENTS_PASS ?= toll
EVENTS_HOST ?= 127.0.0.1
EVENTS_PORT ?= 5432

# Prefix in database name allows multiple preview databases on same env
EVENTS_PREFIX ?=
COMMAND       ?= up

EVENTS_DSN := postgres://$(EVENTS_USER):$(EVENTS_PASS)@$(EVENTS_HOST):$(EVENTS_PORT)

.DEFAULT_GOAL := help
.PHONY: help setup generate migrate/events migrate/events/%

help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

setup: migrate/events ## Setup local development environment

migrate/events: $(addprefix migrate/events/,$(EVENTS_NAME)) ## Run postgres migration schema for configured instance
migrate/events/%: ## Run postgres migration schemas for selected instance
	$(eval time=$(shell date -u +"%Y-%m-%d-%H-%M-%S"))
	docker run \
		--rm \
		--volume "$(CURDIR)/events/migrations":/migrations \
		--network host \
		migrate/migrate:v4.15.2 \
			--path /migrations \
			--database '$(EVENTS_DSN)/$(EVENTS_PREFIX)$*?sslmode=disable' \
			$(COMMAND)
