# gRPC Pet Server

A simple Go server to demonstrate gRPC concepts

## Development

### Prerequisites

- [Buf CLI](https://buf.build/docs/installation)
- Go 1.21 or later (for Go code generation)
- Node.js 20 or later (for TypeScript code generation and npm publishing)

### Linting

```bash
buf lint
```

### Code Generation for Go

```bash
buf generate --template buf.gen.go.yaml
```

### Code Generation for TypeScript

```bash
buf generate --template buf.gen.ts.yaml
```

### Starting the server

```bash
go mod tidy
go run server/main.go
```

### Calling GetStatus

```bash
buf curl \
--schema . \
http://localhost:8080/pet.v1.PetService/GetStatus

# Random response:
# {
#   "status": "PET_STATUS_TRAINING"
# }
```

### Calling SubscribeHeartbeat

```bash
buf curl \
--schema . \
http://localhost:8080/pet.v1.PetService/SubscribeHeartbeat

# Response:
# {
#   "timestampMs": "1766547064733"
# }
# {
#   "timestampMs": "1766547065733"
# }
# ...
```
