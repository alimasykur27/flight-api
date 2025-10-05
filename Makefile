# Build settings
BINARY_DIR=bin
BINARY_NAME=flight-api
MAIN_PATH=./cmd/server/main.go

# Default target
all: build

# Build the application
build:
	mkdir -p $(BINARY_DIR) && \
	go build -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH) && \
	echo "Server binary built in ./bin:" && \
	ls -lh $(BINARY_DIR)

# Run the application
run:
	go run $(MAIN_PATH)

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_DIR}/$(BINARY_NAME)
	rm -f ${BINARY_DIR}/coverage.out

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=${BINARY_DIR}/coverage.out
	go tool cover -html=${BINARY_DIR}/coverage.out

# Build Docker image
docker-build:
	docker build -t ${BINARY_NAME}

# Run Docker Image
docker-run:
	docker run -p 3000:3000 --env-file .env ${BINARY_NAME}\\

# Create DB Test
docker-postgres-test:
	docker exec -it flight_api_db psql -U postgres -d flightdb -c "CREATE DATABASE flightdb_test;"

# Install dependencies
deps:
	go mod tidy
	go mod verify

# Start development environment
dev:
	docker-compose up -d

# Stop development environment
down:
	docker-compose down

# Create a new migration
migrate-create:
	@read -p "Enter migration name : " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Build migration tool
build-migrate:
	go build -o migrate ./cmd/migrate/main.go

# Run database migrations
migrate-up:
	go run ./cmd/migrate/main.go --up

# Rollback migrations
migrate-down:
	go run ./cmd/migrate/main.go --down

# View migration status
migrate-status:
	go run ./cmd/migrate/main.go --status

.PHONY: 
	all \
	build \
	run \
	clean \
	test \
	test-coverage \
	docker-build \
	docker-run \
	deps \
	dev \
	down \
	migrate-create \
	migrate-up \
	migrate-down \
	migrate-status