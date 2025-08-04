# CRUSH.md - Development Guide

## Build/Test Commands
- **Build**: `go build -o rara-membership .`
- **Test all**: `go test ./...`
- **Test single package**: `go test ./path/to/package`
- **Test with coverage**: `go test -cover ./...`
- **Generate Ent code**: `go generate ./ent`
- **Format code**: `go fmt ./...`
- **Vet code**: `go vet ./...`
- **Run CLI**: `./go run main.go --help`
- **List members**: `go run main.go member list`
- **Create member**: `./go run main.go member create`
- **Start server**: `./go run main.go serve`

## Code Style Guidelines

### Imports
- Standard library imports first, then third-party, then local packages
- Use `goimports` for automatic import management

### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for exported types and functions
- Use ALL_CAPS for constants
- Package names should be lowercase, single word when possible

### Types & Error Handling
- Always handle errors explicitly, never ignore them
- Use `fmt.Errorf` for wrapping errors with context
- Define custom error types for domain-specific errors
- Use pointer receivers for methods that modify the receiver

### Ent Framework
- Create new database models with `go run -mod=mod entgo.io/ent/cmd/ent new <ModelName>`; this will place the schema definitions in `ent/schema/`
- Run `go generate ./ent` after schema changes
- Follow Ent naming conventions for fields and edges
- Run `go run -mod=mod entgo.io/ent/cmd/ent describe ent/schema` to describe the models
- Database file: `membership.db` (auto-created, in .gitignore)

### General
- Keep functions small and focused
- Use meaningful variable names
- Add godoc comments for exported functions and types