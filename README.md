# Config Server

A configuration management server with gRPC and HTTP transport layers that can store configuration files using MongoDB GridFS or local filesystem.

## Prerequisites

- Go 1.23 or later
- Protocol Buffers compiler (protoc)
- MongoDB (if not using file-based storage)

## Setup

1. Install dependencies:
```
make deps
```

2. Generate code from proto files:
```
make proto
```

3. Build the server:
```
make build
```

## Configuration

Edit `config.yaml` to configure the server:

```yaml
mongoURI: "mongodb://localhost:27017"
bind:
  http: 8080
  grpc: 50051
use_file: false  # Set to true to use file-based storage instead of MongoDB
```

## Running

Start the server:
```
make run
```

Or run with custom config:
```
./config_server --cfg=custom_config.yaml
```

## Testing

Run tests:
```
make test
```

## Clean

Remove build artifacts:
```
make clean
```