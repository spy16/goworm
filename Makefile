all: clean fmt test build

fmt:
	@echo "Formatting..."
	@goimports -l -w ./

install:
	@echo "Installing goworm to GOBIN..."
	@go install

clean:
	@echo "Cleaning up..."
	@rm -rf ./bin
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -cover ./...

test-verbose:
	@echo "Running tests..."
	@go test -v -cover ./...

build:
	@echo "Building..."
	@mkdir -p ./bin
	@go build -o ./bin/goworm 
