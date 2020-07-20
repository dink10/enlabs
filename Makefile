.PHONY: default
default: up

.PHONY: lint
lint: ## Run linter for all project files
	@echo "Running linter..."
	@golangci-lint run
	@echo "Done"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	bash -c "go test tests/main_test.go -count=1 -v"
	@echo "Done"

.PHONY: swag
swag: ## Run generating swagger documentation
	bash -c "swag init -g cmd/api/main.go"

.PHONY: build
build: ## Run build project
	bash -c "docker-compose -f deployments/docker-compose.yml build"

.PHONY: up
up: ## Run build and up project
	bash -c "docker-compose -f deployments/docker-compose.yml up --build"

.PHONY: start
start: ## Start project
	bash -c "docker-compose -f deployments/docker-compose.yml start"

.PHONY: stop
stop: ## Stop project
	bash -c "docker-compose -f deployments/docker-compose.yml stop"
