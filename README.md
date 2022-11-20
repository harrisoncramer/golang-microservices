# Go Microservices

This repository holds a series of microservices build with Go.

Please refer to the makefile for instructions on how to start up the system and tear it down.

# Commands

## `make build` 

Builds all of the microservices.

## `make up`

Starts all of the microservices.

## `make down`

Stops all microservices.

## `make up_build`

Builds and starts all microservices.

## `make debug SERVICE=authentication`

Rebuilds and restarts all the microservices, including the one specified with the DEBUG argument, which is exposed on port `9080` for a debugger to attach to.

All microservices can restart in debug mode with this command.

## `make rebuild SERVICE=authentication`

Rebuilds a specific microservice.

## `make logs`

Show logs from all the microservices. Wrapper around docker compose.

## `make ps`

Show the currently running microservices. Wrapper around docker compose.

## `make vite`

Starts up the frontend.

## `make vite-stop`

Stops the frontend.
