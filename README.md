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
