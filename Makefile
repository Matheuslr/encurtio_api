# ─────────────────────────────────────────────────────────
# Encurtio – URL Shortener
# ─────────────────────────────────────────────────────────

APP_NAME   := encurtio
BINARY     := ./bin/$(APP_NAME)
CMD        := ./cmd/api
GOFLAGS    := -v

# ─────── Build ───────

.PHONY: build
build: ## Build the binary
	@echo "🔨 Building $(APP_NAME)..."
	go build $(GOFLAGS) -o $(BINARY) $(CMD)

.PHONY: run
run: ## Run locally (requires Cassandra running)
	go run $(CMD)

.PHONY: air
air: ## Run with hot-reload (requires air installed)
	air

# ─────── Test ───────

.PHONY: test
test: ## Run all tests
	go test -v -count=1 ./test/...

.PHONY: test-cover
test-cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./test/...
	go tool cover -func=coverage.out
	@echo "📄 HTML report: go tool cover -html=coverage.out"

.PHONY: test-cover-html
test-cover-html: test-cover ## Open coverage in browser
	go tool cover -html=coverage.out

# ─────── Lint / Format ───────

.PHONY: fmt
fmt: ## Format code
	gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: vet fmt ## Lint (vet + fmt)

# ─────── Dependencies ───────

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: deps
deps: tidy ## Download dependencies
	go mod download

# ─────── Docker ───────

.PHONY: docker-up
docker-up: ## Start all containers (Cassandra + API)
	docker compose up -d --build

.PHONY: docker-down
docker-down: ## Stop all containers
	docker compose down

.PHONY: docker-logs
docker-logs: ## Tail container logs
	docker compose logs -f

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Restart all containers

# ─────── Database ───────

.PHONY: cassandra-up
cassandra-up: ## Start only Cassandra
	docker compose up -d cassandra

.PHONY: cassandra-cql
cassandra-cql: ## Open cqlsh shell
	docker exec -it encurtio-cassandra cqlsh

# ─────── Clean ───────

.PHONY: clean
clean: ## Remove build artefacts
	rm -rf bin/ coverage.out

# ─────── Help ───────

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
