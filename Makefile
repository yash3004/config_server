.PHONY: all proto build test run clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=config_server
MAIN_PATH=./cmd/server

# Proto parameters
PROTOC=protoc
PROTO_DIR=./proto
PROTO_OUT=./proto
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

all: proto build

# Install dependencies
deps:
	$(GOGET) -u google.golang.org/protobuf/cmd/protoc-gen-go
	$(GOGET) -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	$(GOMOD) tidy

# Generate code from proto files
proto:
	$(PROTOC) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	$(GOTEST) -v ./...

# Run the server
run:
	./$(BINARY_NAME) --cfg=config.yaml

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f $(PROTO_OUT)/*.pb.go