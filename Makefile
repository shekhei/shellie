.PHONY: all clean proto build

# Binary names and paths
CLIENT_BINARY = bin/client
SERVICE_BINARY = bin/service
PROTO_DIR = pb
PB_DIR = pb

# Go source files
CLIENT_SRC = client/main.go
SERVICE_SRC = service/main.go

all: proto build

# Create necessary directories
$(shell mkdir -p bin pb)

# Compile protocol buffers
proto:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Build both binaries
build: $(CLIENT_BINARY) $(SERVICE_BINARY)

# Build client binary
$(CLIENT_BINARY): $(CLIENT_SRC) proto
	go build -o $(CLIENT_BINARY) $(CLIENT_SRC)

# Build service binary
$(SERVICE_BINARY): $(SERVICE_SRC) proto
	go build -o $(SERVICE_BINARY) $(SERVICE_SRC)

# Clean built files
clean:
	rm -rf $(CLIENT_BINARY) $(SERVICE_BINARY) $(PB_DIR)/*.pb.go
