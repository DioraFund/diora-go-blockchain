# Diora Blockchain Makefile
.PHONY: help build clean run dev test docker-stop docker-start docker-logs install-deps check-status stop

# Default target
help:
	@echo "Diora Blockchain Development Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build all components"
	@echo "  make run          - Run all services locally"
	@echo "  make dev          - Start development mode"
	@echo "  make docker-start  - Start with Docker Compose"
	@echo "  make docker-stop   - Stop Docker Compose services"
	@echo "  make docker-logs   - Show Docker logs"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make install-deps - Install dependencies"
	@echo "  make check-status - Check service status"
	@echo "  make stop         - Stop all services"

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Build all components
build:
	@echo "$(BLUE)Building Diora Blockchain...$(NC)"
	@echo "$(YELLOW)Building blockchain node...$(NC)"
	cd core && go build -o ../build/diora .
	@echo "$(YELLOW)Building API server...$(NC)"
	cd api && go build -o ../build/api-server .
	@echo "$(GREEN)✓ Build completed$(NC)"

# Install dependencies for all components
install-deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	go mod tidy
	go mod download
	@echo "$(YELLOW)Installing Node.js dependencies for web...$(NC)"
	cd web && npm install
	@echo "$(YELLOW)Installing Node.js dependencies for explorer...$(NC)"
	cd explorer && npm install
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

# Run all services locally
run: build
	@echo "$(BLUE)Starting Diora Blockchain services...$(NC)"
	@echo "$(YELLOW)Starting blockchain node...$(NC)"
	./build/diora node start --datadir=./data --network=testnet &
	@echo "$(YELLOW)Starting API server...$(NC)"
	./build/api-server --blockchain-rpc=http://localhost:8545 &
	@echo "$(YELLOW)Starting web interface...$(NC)"
	cd web && npm run dev &
	@echo "$(YELLOW)Starting blockchain explorer...$(NC)"
	cd explorer && npm run dev &
	@echo "$(GREEN)✓ All services started$(NC)"
	@echo "$(BLUE)Access URLs:$(NC)"
	@echo "  Blockchain Node: http://localhost:8545"
	@echo "  API Server: http://localhost:8080"
	@echo "  Web Interface: http://localhost:3000"
	@echo "  Blockchain Explorer: http://localhost:3001"
	@echo ""
	@echo "$(YELLOW)Press Ctrl+C to stop all services$(NC)"
	@wait

# Development mode with hot reload
dev: install-deps
	@echo "$(BLUE)Starting development mode...$(NC)"
	@echo "$(YELLOW)Starting blockchain node in dev mode...$(NC)"
	./build/diora node start --dev --datadir=./data &
	sleep 2
	@echo "$(YELLOW)Starting API server in dev mode...$(NC)"
	cd api && go run server.go --debug &
	sleep 2
	@echo "$(YELLOW)Starting web interface with hot reload...$(NC)"
	cd web && npm run dev &
	sleep 2
	@echo "$(YELLOW)Starting explorer with hot reload...$(NC)"
	cd explorer && npm run dev &
	@echo "$(GREEN)✓ Development mode started$(NC)"
	@echo "$(BLUE)Access URLs:$(NC)"
	@echo "  Blockchain Node: http://localhost:8545"
	@echo "  API Server: http://localhost:8080"
	@echo "  Web Interface: http://localhost:3000"
	@echo "  Blockchain Explorer: http://localhost:3001"
	@echo ""
	@echo "$(YELLOW)Development mode with hot reload enabled$(NC)"
	@wait

# Docker Compose commands
docker-start:
	@echo "$(BLUE)Starting Diora with Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Services started with Docker$(NC)"
	@echo "$(BLUE)Access URLs:$(NC)"
	@echo "  Blockchain Node: http://localhost:8545"
	@echo "  API Server: http://localhost:8080"
	@echo "  Web Interface: http://localhost:3000"
	@echo "  Blockchain Explorer: http://localhost:3001"
	@echo "  Monitoring: http://localhost:9090 (Prometheus)"
	@echo "  Dashboard: http://localhost:3002 (Grafana)"

docker-stop:
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

docker-logs:
	@echo "$(BLUE)Showing Docker logs...$(NC)"
	docker-compose logs -f

docker-rebuild:
	@echo "$(BLUE)Rebuilding Docker images...$(NC)"
	docker-compose build --no-cache
	docker-compose up -d

# Test commands
test:
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./tests/...
	cd web && npm test
	cd explorer && npm test

test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf build/
	rm -rf data/
	rm -rf coverage.out
	rm -rf coverage.html
	go clean -cache
	cd web && rm -rf node_modules/ .next/
	cd explorer && rm -rf node_modules/ .next/
	@echo "$(GREEN)✓ Clean completed$(NC)"

# Stop all running services
stop:
	@echo "$(RED)Stopping all Diora services...$(NC)"
	@echo "$(YELLOW)Stopping local services...$(NC)"
	pkill -f "diora" || true
	pkill -f "api-server" || true
	pkill -f "npm run dev" || true
	pkill -f "node" || true
	@echo "$(GREEN)✓ All services stopped$(NC)"

# Check status of services
check-status:
	@echo "$(BLUE)Checking service status...$(NC)"
	@echo "$(YELLOW)Blockchain Node:$(NC) $$(shell curl -s http://localhost:8545/health > /dev/null && echo "$(GREEN)Running$(NC)" || echo "$(RED)Stopped$(NC))"
	@echo "$(YELLOW)API Server:$(NC) $$(shell curl -s http://localhost:8080/health > /dev/null && echo "$(GREEN)Running$(NC)" || echo "$(RED)Stopped$(NC))"
	@echo "$(YELLOW)Web Interface:$(NC) $$(shell curl -s http://localhost:3000 > /dev/null && echo "$(GREEN)Running$(NC)" || echo "$(RED)Stopped$(NC))"
	@echo "$(YELLOW)Explorer:$(NC) $$(shell curl -s http://localhost:3001 > /dev/null && echo "$(GREEN)Running$(NC)" || echo "$(RED)Stopped$(NC))"

# Quick development setup
setup: install-deps build
	@echo "$(GREEN)✓ Development environment setup complete$(NC)"
	@echo "$(BLUE)Run 'make dev' to start development mode$(NC)"

# Production build
prod-build:
	@echo "$(BLUE)Building for production...$(NC)"
	@echo "$(YELLOW)Building blockchain node...$(NC)"
	cd core && go build -ldflags="-s -w" -o ../build/diora-linux .
	@echo "$(YELLOW)Building web interface...$(NC)"
	cd web && npm run build
	@echo "$(YELLOW)Building explorer...$(NC)"
	cd explorer && npm run build
	@echo "$(GREEN)✓ Production build completed$(NC)"

# Generate configuration files
config:
	@echo "$(BLUE)Generating configuration files...$(NC)"
	@echo "$(YELLOW)Creating config directory...$(NC)"
	mkdir -p config
	@echo "$(YELLOW)Generating blockchain config...$(NC)"
	echo '[network]' > config/config.toml
	echo 'chain_id = 1337' >> config/config.toml
	echo 'network_id = 1' >> config/config.toml
	echo 'block_time = 6' >> config/config.toml
	echo 'gas_limit = 30000000' >> config/config.toml
	echo 'min_gas_price = 1000000000' >> config/config.toml
	echo '' >> config/config.toml
	echo '[consensus]' >> config/config.toml
	echo 'stake_amount = 1000000' >> config/config.toml
	echo 'validator_count = 100' >> config/config.toml
	echo 'unbonding_period = 604800' >> config/config.toml
	echo '' >> config/config.toml
	echo '[api]' >> config/config.toml
	echo 'rpc_port = 8545' >> config/config.toml
	echo 'ws_port = 8546' >> config/config.toml
	echo 'max_connections = 1000' >> config/config.toml
	echo 'rate_limit = 100' >> config/config.toml
	echo '' >> config/config.toml
	echo '[storage]' >> config/config.toml
	echo 'data_dir = "./data"' >> config/config.toml
	echo 'db_type = "leveldb"' >> config/config.toml
	echo 'cache_size = 1024' >> config/config.toml
	echo '' >> config/config.toml
	echo '[logging]' >> config/config.toml
	echo 'level = "info"' >> config/config.toml
	echo 'file = "./logs/diora.log"' >> config/config.toml
	echo 'max_size = 100' >> config/config.toml
	@echo "$(GREEN)✓ Configuration files generated$(NC)"

# Create necessary directories
init-dirs:
	@echo "$(BLUE)Creating necessary directories...$(NC)"
	mkdir -p data logs config build
	@echo "$(GREEN)✓ Directories created$(NC)"

# Quick start command
quickstart: init-dirs config build
	@echo "$(GREEN)✓ Quickstart complete!$(NC)"
	@echo "$(BLUE)Run 'make run' to start all services$(NC)"
