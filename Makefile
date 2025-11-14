ifeq ($(OS),Windows_NT)
    BIN_EXT := .exe
else
    BIN_EXT :=
endif

BIN_DIR := bin
API_DIR := api

GOLANGCI_LINT_BIN := $(BIN_DIR)/golangci-lint$(BIN_EXT)
OAPI_CODEGEN_BIN := $(BIN_DIR)/oapi-codegen$(BIN_EXT)
GOFUMPT_BIN := $(BIN_DIR)/gofumpt$(BIN_EXT)
GOOSE_BIN := $(BIN_DIR)/goose$(BIN_EXT)

API_SCHEMA := openapi.yaml
CODEGEN_CONFIG := codegen.yaml

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

GOOSE_DRIVER = postgres
GOOSE_DBSTRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)
MIGRATIONS_DIR = migrations

.PHONY: all dev run devenv-start debug test deps tools clean codegen fmt \
        migrate-up migrate-down migrate-status postgres-up postgres-stop postgres-health connect-db

dev: devenv-start run

all: codegen lint fmt

run:
	@go run service/cmd/service

devenv-start:
	@cp -f .env.example .env

debug:
	@docker compose --env-file .env up -d

postgres-stop:
	@docker compose --env-file .env stop

postgres-health:
	@docker exec my-postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}

deps:
	@go mod download
	@go mod tidy
	@go mod verify

tools: $(OAPI_CODEGEN_BIN) $(GOFUMPT_BIN) $(GOOSE_BIN) $(GOLANGCI_LINT_BIN)

codegen: $(OAPI_CODEGEN_BIN)
	@echo "Generating code from API schema..."
	@$(OAPI_CODEGEN_BIN) -config $(API_DIR)/$(CODEGEN_CONFIG) $(API_DIR)/$(API_SCHEMA)

$(OAPI_CODEGEN_BIN):
	@echo "building oapi-codegen"
	@go build -o $(OAPI_CODEGEN_BIN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	@$(OAPI_CODEGEN_BIN) --version

fmt: $(GOFUMPT_BIN)
	@$(GOFUMPT_BIN) -l -w .

$(GOFUMPT_BIN):
	@echo "building gofumpt"
	@go build -o $(GOFUMPT_BIN) mvdan.cc/gofumpt
	@$(GOFUMPT_BIN) --version

migrate-up: $(GOOSE_BIN)
	@echo "Applying migrations from $(MIGRATIONS_DIR)..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) up

migrate-down: $(GOOSE_BIN)
	@echo "Rolling back migrations..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) down

migrate-status: $(GOOSE_BIN)
	@echo "Migration status in $(MIGRATIONS_DIR):"
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) status

$(GOOSE_BIN):
	@echo "building goose"
	@go build -o $(GOOSE_BIN) github.com/pressly/goose/v3/cmd/goose
	@$(GOOSE_BIN) --version

lint: $(GOLANGCI_LINT_BIN)
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT_BIN) run

$(GOLANGCI_LINT_BIN):
	@echo "Installing golangci-lint..."
	go build -o $(GOLANGCI_LINT_BIN) github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GOLANGCI_LINT_BIN) --version

connect-db:
	@psql $(POSTGRES_CONNSTRING)