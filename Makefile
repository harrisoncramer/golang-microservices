SHELL := /bin/bash

FRONT_END_BINARY=frontApp

BROKER_BINARY=brokerApp.bin
BROKER_DEBUG_BINARY=brokerDebug.bin

AUTH_BINARY=authApp.bin
AUTH_DEBUG_BINARY=authDebug.bin

LOGGER_BINARY=loggerApp.bin
LOGGER_DEBUG_BINARY=loggerDebug.bin

# up: Starts all containers in the background without forcing a build
up:
	@echo "Starting docker images"
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "Docker images started"

# up_build: stops docker-compose (if running) and builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Building (when required) and starting docker images"
	@docker-compose -f docker/docker-compose.yml up --build -d
	@echo "Docker images built and started!"

# Passing an environment variable with DEBUG='authentication' will cause the authentication-service to be built with
# the correct debugger flags and run via delve. The DEBUG value should always be the first part of the service's name,
# prior to the `-service` suffix. Delve will the be exposed on port 9080 locally for connection.

# debug: stops docker-compose (if running) and builds project specified by DEBUG value with docker-compose
debug: debug_check build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Building and starting docker images (with ${DEBUG} in debug mode)"
	@docker-compose -f docker/docker-compose.yml -f docker/docker-compose.debug.$$DEBUG.yml up --build -d
	@echo "Debug mode started for ${DEBUG} service! Connect on port 9080."

debug_check:
	 @[ "${DEBUG}" ] && echo "Debugging ${DEBUG}..." || ( echo "DEBUG value is not set!"; exit 1 )

# down: stop docker compose
down:
	@echo "Stopping docker compose"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Done!"

# build_logger: builds the logger binary as a linux executable
build_logger: 
	@echo "Building logger binary..."
	@if [ "${DEBUG}" = 'logger' ]; then \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${LOGGER_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api; \
	fi

# build_broker: builds the broker binary as a linux executable
build_broker: 
	@echo "Building broker binary..."
	@if [ "${DEBUG}" = 'broker' ]; then \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${BROKER_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api; \
	fi

# build_auth: builds the auth binary as a linux executable
build_auth: 
	@echo "Building auth binary..."
	@if [ "${DEBUG}" = 'authentication' ]; then \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ${AUTH_DEBUG_BINARY} ./cmd/api; \
	else \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api; \
	fi

# logs: shows logs from the containers
logs:
	@docker-compose -f docker/docker-compose.yml logs

# ps: shows currently running docker processes
ps:
	@docker-compose -f docker/docker-compose.yml ps

# start: Starts up the Vite frontend
start:
	@echo "Starting front end"
	@cd ./frontend && npm run dev

# stop: Stops the Vite frontend
stop:
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "node .*/golang-microservices/frontend/node_modules/.bin/vite"
