.PHONY: run swag clean install setup

# Go commands
GO=go
SWAG=swag

# Binary name
BINARY_NAME=server

# Run the application
run: swag
	$(GO) run cmd/server/main.go

# Generate Swagger documentation
swag:
	$(SWAG) init -g cmd/server/main.go

# Install dependencies
install:
	$(GO) mod tidy
	$(GO) get -u github.com/swaggo/swag/cmd/swag
	$(GO) get -u github.com/swaggo/gin-swagger
	$(GO) get -u github.com/swaggo/files
	$(GO) get -u github.com/golang-jwt/jwt/v5

# Clean build files and generated documentation
clean:
	rm -rf docs
	rm -f $(BINARY_NAME)

# Setup: install dependencies and generate Swagger docs
setup: install swag 