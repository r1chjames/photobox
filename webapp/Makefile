help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean working directory
	rm -f .env
	rm -rf build/*
	rm -rf node_modules

.PHONY: install
install: ## Install NPM dependencies
	npm install

.PHONY: run
run: ## Run web app
	npm run start

.PHONY: test
test: install ## Run web app tests
	npm test

.PHONY: build_app
build_app: clean install ## Build web app
	npm run build

.PHONY: build_app_artifact
build_app_artifact: build_app ## Package web app artifact
	tar -zcvf webapp.tar.gz build

build_all: build_app_artifact
