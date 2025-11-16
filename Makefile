ifeq ($(OS),Windows_NT)
	BIN_EXT := .exe
else
	BIN_EXT :=
endif

ifneq (,$(wildcard .env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif
BIN_DIR := bin
API_DIR := api

GOLANGCI_LINT_BIN := $(BIN_DIR)/golangci-lint$(BIN_EXT)
OAPI_CODEGEN_BIN := $(BIN_DIR)/oapi-codegen$(BIN_EXT)
GOFUMPT_BIN := $(BIN_DIR)/gofumpt$(BIN_EXT)
GOOSE_BIN := $(BIN_DIR)/goose$(BIN_EXT)

API_SCHEMA := openapi.yaml
CODEGEN_CONFIG := codegen.yaml

MIGRATIONS_DIR := migrations
GOOSE_DRIVER := postgres

GOOSE_DBSTRING := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

.PHONY: all devenv-start tools codegen fmt lint migrate-up migrate-down migrate-status

all: tidy codegen lint fmt devenv-start

tidy:
	go mod tidy

devenv-start:
	@cp -f .env.example .env

run: devenv-start
	@docker compose --env-file .env up -d

tools: $(OAPI_CODEGEN_BIN) $(GOFUMPT_BIN) $(GOOSE_BIN) $(GOLANGCI_LINT_BIN)

$(OAPI_CODEGEN_BIN):
	@go build -o $(OAPI_CODEGEN_BIN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	@$(OAPI_CODEGEN_BIN) --version

$(GOFUMPT_BIN):
	@go build -o $(GOFUMPT_BIN) mvdan.cc/gofumpt
	@$(GOFUMPT_BIN) --version

$(GOOSE_BIN):
	@go build -o $(GOOSE_BIN) github.com/pressly/goose/v3/cmd/goose
	@$(GOOSE_BIN) --version

$(GOLANGCI_LINT_BIN):
	@go build -o $(GOLANGCI_LINT_BIN) github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GOLANGCI_LINT_BIN) --version

codegen: $(OAPI_CODEGEN_BIN)
	@$(OAPI_CODEGEN_BIN) -config $(API_DIR)/$(CODEGEN_CONFIG) $(API_DIR)/$(API_SCHEMA)

fmt: $(GOFUMPT_BIN)
	@$(GOFUMPT_BIN) -l -w .

lint: $(GOLANGCI_LINT_BIN)
	@$(GOLANGCI_LINT_BIN) run

migrate-up: $(GOOSE_BIN)
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) up

migrate-down: $(GOOSE_BIN)
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) down

migrate-status: $(GOOSE_BIN)
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_BIN) -dir $(MIGRATIONS_DIR) status
