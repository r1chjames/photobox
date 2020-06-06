help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean working directory
	rm -f .env
	rm -rf build/*

.PHONY: install
install: ## Install Go dependencies
	npm install

.PHONY: run
run:
	npm run start

.PHONY: test
test: install## Start Kafka and Kafka Connect containers and run test hardness pack
	npm test

.PHONY: build_app
build_app: clean install
	npm run build

.PHONY: build_app_artifact
build_app_artifact: build_app
	tar -zcvf webapp.tar.gz build

build_all: build_app_artifact
