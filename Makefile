SHELL := /bin/bash

# up: Starts all containers in the background without forcing a build
up:
	@echo "Stopping docker images (if running)"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Starting docker images"
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "Docker images started"

# down: Stops all microservices
down:
	@echo "Stopping docker compose"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Done!"

# up_build: stops docker-compose (if running) and builds all projects and starts docker compose
up_build: build
	@echo "Stopping docker images (if running)"
	@docker-compose -f docker/docker-compose.yml down
	@echo "Starting docker images"
	@docker-compose -f docker/docker-compose.yml up -d
	@echo "Docker images started!"

# build: builds all microservices
build: build_broker build_auth build_logger
	@echo "Building (when required) and starting docker images"
	@docker-compose -f docker/docker-compose.yml build
	@echo "Docker images built!"

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

rebuild: rebuild_check
	@docker-compose -f docker/docker-compose.yml build $$TARGET-service

# rebuild_check: Checkes whether, when running `make rebuild`, the TARGET value has actually been set
rebuild_check:
	 @[ "${TARGET}" ] && echo "Rebuilding ${TARGET}..." || ( echo "TARGET value is not set!"; exit 1 )

# debug_check: Checks whether, when running `make debug`, the DEBUG value has actually been set
debug_check:
	 @[ "${DEBUG}" ] && echo "Debugging ${DEBUG}..." || ( echo "DEBUG value is not set!"; exit 1 )

# build_logger: builds the logger binary as a linux executable
build_logger: 
	@echo "Building logger binary..."
	@if [ "${DEBUG}" = 'logger' ]; then \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o 'logger.debug.bin' ./cmd/api; \
	else \
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o 'logger.bin' ./cmd/api; \
	fi

# build_broker: builds the broker binary as a linux executable
build_broker: 
	@echo "Building broker binary..."
	@if [ "${DEBUG}" = 'broker' ]; then \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o 'broker.debug.bin' ./cmd/api; \
	else \
		cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o 'broker.bin' ./cmd/api; \
	fi

# build_auth: builds the auth binary as a linux executable
build_auth: 
	@echo "Building auth binary..."
	@if [ "${DEBUG}" = 'authentication' ]; then \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o 'authentication.debug.bin' ./cmd/api; \
	else \
		cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o 'authentication.bin' ./cmd/api; \
	fi

# logs: shows logs from the containers
logs:
	@docker-compose -f docker/docker-compose.yml logs

# ps: shows currently running docker processes
ps:
	@docker-compose -f docker/docker-compose.yml ps

# start: Starts up the Vite frontend
vite:
	@echo "Starting front end"
	@cd ./frontend && npm run dev

# stop: Stops the Vite frontend
vite-stop:
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "node .*/golang-microservices/frontend/node_modules/.bin/vite"
