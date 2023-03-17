# Makefile for car-pooling-challenge
# vim: set ft=make ts=8 noet
# Copyright Cabify.com
# Licence MIT

# Variables
# UNAME		:= $(shell uname -s)

.EXPORT_ALL_VARIABLES:

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help: ### this screen
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Targets
#
.PHONY: debug
debug:	### Debug Makefile itself
	@echo $(UNAME)

.PHONY: build
build: ### Builds go application and puts generated binary in /target/bin/carpool
	CGO_ENABLED=0 go build -a -o target/bin/carpool ./cmd/carpool/main.go

.PHONY: go-run
go-run: build ### Builds go application and runs it. The necessary infrastructure will not be raised with this command. To pull up the necessary dependencies use 'make run' instead.
	target/bin/carpool

.PHONY: dockerize
docker: ### Builds go application using Dockerfile
	docker build -t car-pooling-challenge:latest .

.PHONY: test.acceptance
test.acceptance: docker ### Runs challenge acceptance tests using 'harness' docker-compose service
	CABIFY_CHALLENGE_TESTCASE=acceptance docker-compose down && docker-compose rm -f && docker-compose up harness --always-recreate-deps --force-recreate --build --abort-on-container-exit

.PHONY: test
test: test.acceptance ### Alias of test.acceptance

.PHONY: run
run: docker ### Runs pooling service using docker-compose
	docker-compose down && docker-compose rm -f && docker-compose up --build pooling

.PHONY: env
env: ### Configures required environment variables to start the application
	@if [ -f .env ]; then echo "The .env file has already been created"; else cp .example.env .env && echo "Copied .example.env to .env"; fi
