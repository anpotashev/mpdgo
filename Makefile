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

# none

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test/coverage: run test coverage
#.PHONY: test/coverage
test/coverage:
	go test -v -cover  -coverprofile=cover.txt ./...

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
# DEBUG
# ==================================================================================== #

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
