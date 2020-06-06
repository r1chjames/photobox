
.PHONY: build-api
build-api:
	cd api && make docker-build


.PHONY: build-webapp
build-webapp:
	cd webapp && make docker-build

.PHONY: build-all 
build-all: build-api build-webapp