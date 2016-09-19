all: clean compile

clean: ## Clean up all the resources
	rm -rf target/*
	go clean -r

compile: ## Compile the project to generate the binary in the target folder
	go get
	go build -ldflags "-X main.version=${VERSION} -X main.minversion=`date -u +.%Y%m%d.%H%M%S` -X main.buildTime=`date  +'%Y-%m-%d'`" -o target/logstash_prom_exporter

test: ## Run the test cases in random order via ginkgo
	# The sleep is there to allow other databases & linked services to start up.
	# This is required by docker-compose. Please bear with this in dev.
	sleep 5s
	go get
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race

run: ## Run the binary
	# The sleep is there to allow other databases & linked services to start up.
	# This is required by docker-compose. Please bear with this in dev.
	sleep 5s
	./target/logstash_prom_exporter

debug: ##create debug file and debugs it woth gdb
	go build -gcflags "-N -l" -o gdb_sandbox
	gdb gdb_sandbox

.PHONY: help

help: ## You can always run this command to see what options are available to you while running the make command
		@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help