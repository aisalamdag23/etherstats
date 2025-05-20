.PHONY: config

# runtime options
COMMIT_HASH = $(shell git rev-parse --short HEAD || echo "$(build_commit)")
TAG         = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match || echo "$(build_tag)")

GO			?= go
DC			?= docker-compose

DBMATE ?= dbmate

DIRENV ?= direnv

LDFLAGS := -X "main.Tag=$(TAG)" \
		   -X "main.CommitHash=$(COMMIT_HASH)"
DOCKER_COMPOSE_CONFIG := ./docker-compose.yml

# Use your ETHERSTATS_POSTGRES_DSN connection string
PSQL ?= psql -d "$$ETHERSTATS_POSTGRES_DSN" -v ON_ERROR_STOP=1

config:
	@cp -n .config.yml.dist .config.yml 2> /dev/null || true
	@cp -n .example.envrc .envrc 2> /dev/null || true
	$(DIRENV) allow # approve changes in envrc

start:
	SPEC_FILE=./.config.yml $(GO) run -ldflags '$(LDFLAGS)' cmd/server/main.go

db-migrate:
	$(DBMATE) -e ETHERSTATS_POSTGRES_DSN migrate

db-new-migration:
	$(DBMATE) new $(name)

db-down:
	$(DBMATE) down

db-reload:
	@$(MAKE) docker-stop
	@sleep 5
	@$(MAKE) docker-start
	@sleep 7
	@$(MAKE) db-migrate

docker-start:
	$(DC) -f $(DOCKER_COMPOSE_CONFIG) up -d

docker-stop:
	$(DC) -f $(DOCKER_COMPOSE_CONFIG) down
