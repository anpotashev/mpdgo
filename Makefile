include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/gompd: run the cmd/gompd application
.PHONY: run/gompd
run/gompd:
	go run ./cmd/gompd

## run/rest-api: run the cmd/rest-api/main.go application
.PHONY: run/rest-api
run/rest-api:
	go run ./cmd/rest-api/main.go

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

BINARY_NAME=mpdapp
BUILD_DIR=build

PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	darwin/amd64 \
	darwin/a

## build/api: build the cmd/api application
.PHONY: build
build:
	@echo 'Building...'
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		ext=""; \
		if [ "$$GOOS" = "windows" ]; then ext=".exe"; fi; \
		output="$(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH$$ext"; \
		echo " → $$output"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$output ./cmd/ || exit 1; \
	done
#	go build -ldflags='-s' -o=./bin/gompd ./cmd/gompd
#	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

clean:
	@rm -rf $(BUILD_DIR)

# remote debug port
DEBUG_PORT=2345
# path to package
PKG?=./
# test name
TEST?=
## debug-one-test: run dlv for one test
.PHONY: debug-one-test
debug-one-test:
	@if [ -z "$(TEST)" ]; then \
		echo "⚠️  Enter test name: make debug-one-test TEST=TestMyFunc"; \
		exit 1; \
	fi
	dlv test $(PKG) \
		--headless \
		--listen=0.0.0.0:$(DEBUG_PORT) \
		--api-version=2 \
		--accept-multiclient \
		-- -test.run ^$(TEST)$$

## gen/mappers: generate mappers
.PHONY: gen/mappers
gen/mappers:
	go run cmd/mappergenerator/gen_mapper.go
	git add internal/api/dto/generated_mapper.go

## test/coverage: run test coverage
#.PHONY: test/coverage
test/coverage:
	go test -v -cover  -coverprofile=cover.txt ./...