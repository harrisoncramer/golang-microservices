FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp

## up: Starts all containers in the background without forcing a build
up:
	@echo "Starting docker images"
	docker-compose up -d
	@echo "Docker images started"

## up_build: stops docker-compose (if running) and builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	docker-compose down
	@echo "Building (when required) and starting docker images"
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose"
	docker-compose down
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger: 
	@echo "Building logger binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker: 
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth: 
	@echo "Building auth binary..."
	cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## Starts Vite
start:
	@echo "Starting front end"
	cd ./frontend && npm run dev

## Stops Vite
stop:
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "node .*/golang-microservices/frontend/node_modules/.bin/vite"
