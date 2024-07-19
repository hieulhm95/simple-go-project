UNAME := $(shell uname)
PWD = $(shell pwd)

GITHUB_SHA ?= ct-core-standard
COMPANY_REPO ?= github.com/carousell
SRC_ROOT = `go list -f {{.Root}} .`
SRC = `go list -f {{.Dir}} ./... | grep -v /vendor/ `

ifeq ($(shell uname -m), arm64)
	TAGS = dynamic
else
	TAGS =
endif

install:
	@echo "==> Installing tools..."
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/incu6us/goimports-reviser/v3@latest
	@brew install golangci/tap/golangci-lint
	@brew upgrade golangci/tap/golangci-lint

fmt:
	@echo "==> Formatting source code..."
	@echo "==> Running gofumpt..."
	@gofumpt -w $(SRC)
	@echo "==> Running goimports-reviser..."
	@go list -f {{.Dir}} ./... | grep -v /vendor/ | sed 's|$(SRC_ROOT)|.|g' | xargs -I {} goimports-reviser -rm-unused -company-prefixes $(COMPANY_REPO) -file-path {}

lint:
	@echo "==> Running lint check..."
	@golangci-lint --config setup/.golangci.yml run
	@go vet -tags=$(TAGS) `go list -f {{.Dir}} ./... | grep -v /vendor/`
	
test:
	@echo "==> Running test"
	@go clean -testcache ./...
	go test -vet=off ./... -p 1 -race -cover -coverprofile=coverage.out

generate:
	@echo "==> generate mock"
	go generate ./...

dev-up:
	@docker compose \
		-f setup/docker-compose.dev.yml \
		-p $(GITHUB_SHA) up --build -d \
		--remove-orphans

dev-down:
	@docker compose \
		-f setup/docker-compose.dev.yml \
		-p $(GITHUB_SHA) down \
		-v --rmi local

test-up:
	@docker-compose \
		-f setup/docker-compose.test.yml \
		-p $(GITHUB_SHA) up \
		--force-recreate \
		--abort-on-container-exit \
		--exit-code-from app \
		--build

test-down:
	@docker-compose \
		-f setup/docker-compose.test.yml \
		-p $(GITHUB_SHA) down \
		-v --rmi local

report:
	@echo "==> Coverage report"
	go tool cover -func=coverage.out

report-html:
	@echo "==> Coverage report html"
	go tool cover -html=coverage.out

.PHONY: test dev-up dev-down test-up test-down lint fmt install generate report report-html 
