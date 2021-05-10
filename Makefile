### Required tools
GOTOOLS_CHECK = go golangci-lint goverreport goimports
PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)
export GO111MODULE=on

all: check_tools ensure-deps build fmt imports linter test

### Tools & dependencies
check_tools:
	@echo "Found tools: $(foreach tool,$(GOTOOLS_CHECK),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

### Build
build:
	@echo "==> Building..."
	go build -o ./cmd/server/api ./cmd/server

test:
	@echo "==> Running tests..."
	go test ./...

test-race:
	@echo "==> Running tests..."
	go test -race ./...

test-cover:
	@echo "==> Running tests with coverage..."
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	goverreport -coverprofile=/tmp/coverage.out -sort=block -order=desc -threshold=85 || (echo -e "**********Minimum test coverage was not reached(85%)**********")
	go tool cover -html=/tmp/coverage.out

### Formatting, linting, and deps
fmt:
	@echo "==> Running format..."
	go fmt ./...

linter:
	@echo "==> Running linter..."
	golangci-lint run ./...

ensure-deps:
	@echo "=> Syncing dependencies with go mod tidy"
	@go mod tidy

### Run binaries
run:
	@echo "==> Running local api command..."
	go run cmd/server/main.go

imports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all check_tools test test-cover fmt linter ensure-deps imports
