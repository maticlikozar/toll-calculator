DEVBOX_DOMAIN      = toll.test
DEVBOX_PROJECTS    = api
BOOTSTRAP_SERVICES = dnsmasq proxy events

.DEFAULT_GOAL := help
.PHONY: help devbox/setup devbox/start devbox/start/services devbox/stop devbox/restart devbox/dns devbox/containers devbox/container

help: ## Help information
	@echo "------------------------------------------------------------------------"
	@echo " Get started:"
	@echo "  - \033[36mmake devbox/setup\033[0m             # Setup dev environment (run once)"
	@echo "  - \033[36mmake devbox/database\033[0m          # Initial database setup"
	@echo "  - \033[36mmake devbox/start\033[0m             # Run dev environment"
	@echo ""
	@echo "  - \033[36mmake devbox/start/<service>\033[0m   # Run single service [$(DEVBOX_PROJECTS)]"
	@echo "  - \033[36mmake devbox/stop\033[0m              # Stop dev environment"
	@echo "------------------------------------------------------------------------"
	@echo " Target list:"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@printf "\nDocs:\n\n"
	@printf "See docs folder for detailed documentation and guides.\n\n"

devbox/setup: devbox/dns devbox/containers

devbox/start/services: ## Start common services
	docker-compose up -d $(BOOTSTRAP_SERVICES)

devbox/start: devbox/start/services database/migrate ## Starts devbox
	@for project in $(DEVBOX_PROJECTS); do \
		$(MAKE) devbox/start/$$project ; \
	done

	docker-compose start

devbox/start/%: devbox/start/services ## Starts service (e.g. make devbox/start/portal)
	@echo "Starting devbox project '$*'"
	docker-compose up -d $*

devbox/stop: ## Stops and removes services
	docker-compose down

devbox/restart/%: ## Restart service
	docker-compose restart $*

devbox/dns: ## Setup DNS resolving for preview domains
	sudo mkdir -p /etc/resolver
	echo "nameserver 127.0.0.1" | sudo tee "/etc/resolver/$(DEVBOX_DOMAIN)"
	echo "port 12215" | sudo tee -a "/etc/resolver/$(DEVBOX_DOMAIN)"

devbox/container/%: ## Build a single development container
	$(MAKE) -C ./resources/docker build/$*

devbox/containers: ## Build development containers
	$(MAKE) -C ./resources/docker build

devbox/database: ## Initial database setup.
	docker-compose up -d events
	# Wait for Database
	@sleep 30
	$(MAKE) -C ./code/sql setup

database/migrate: ## Run migration schemas and seed scripts
	# Wait for Database
	@sleep 30
	$(MAKE) -C ./code/sql migrate/events
