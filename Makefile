SHELL := /bin/bash

FRONT_END_BINARY=frontApp

BROKER_BINARY=brokerApp
BROKER_DEBUG_BINARY=brokerDebug

AUTH_BINARY=authApp
AUTH_DEBUG_BINARY=authDebug

LOGGER_BINARY=loggerApp
LOGGER_DEBUG_BINARY=loggerDebug

# up: Starts all containers in the background without forcing a build
up:
	@echo "Starting docker images"
	docker-compose -f docker/docker-compose.yml up -d
	@echo "Docker images started"

# up_build: stops docker-compose (if running) and builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	docker-compose -f docker/docker-compose.yml down
	@echo "Building (when required) and starting docker images"
	docker-compose -f docker/docker-compose.yml up --build -d
	@echo "Docker images built and started!"

# Passing an environment variable with DEBUG='broker-service' will cause the broker service to be built with
# the correct debugger flags and run via delve. The port of the exposed application can then be connected
# to with a debugger, see the `docker-compose.debug.yml` file for the correct port of the service.

# debug: stops docker-compose (if running) and builds all projects for debugging and starts docker compose
debug: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	docker-compose -f docker/docker-compose.yml down
	@echo "Building and starting docker images (with $DEBUG in debug mode)"
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.debug.yml up --build -d
	@echo "Debug mode started!"

# down: stop docker compose
down:
	@echo "Stopping docker compose"
	docker-compose -f docker/docker-compose.yml down
	@echo "Done!"

# build_logger: builds the logger binary as a linux executable
build_logger: 
	echo "Building logger binary..."
	if [ "${DEBUG}" = 'auth' ]; then \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${LOGGER_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api; \
	fi
	@echo "Done!"

# build_broker: builds the broker binary as a linux executable
build_broker: 
	@echo "Building broker binary..."
	if [ "${DEBUG}" = 'auth' ]; then \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${BROKER_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api; \
	fi

# build_auth: builds the auth binary as a linux executable
build_auth: 
	@echo "Building auth binary..."
	if [ "${DEBUG}" = 'auth' ]; then \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${AUTH_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api; \
	fi
	@echo "Done!"

# start: Starts up the Vite frontend
start:
	@echo "Starting front end"
	cd ./frontend && npm run dev

# stop: Stops the Vite frontend
stop:
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "node .*/golang-microservices/frontend/node_modules/.bin/vite"
