.PHONY: swag build clean build-linux build-windows build-darwin

# Generate Swagger documentation
swag:
	@echo "Generating Swagger documentation..."
	swag init --parseDependency --parseInternal

# Build for multiple platforms and architectures
build:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -o dist/gin-starter-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -o dist/gin-starter-linux-arm64 .
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -o dist/gin-starter-windows-amd64.exe .
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -o dist/gin-starter-darwin-amd64 .
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o dist/gin-starter-darwin-arm64 .
	@echo "Build complete! Binaries are in ./dist"

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p dist
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -o dist/gin-starter-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -o dist/gin-starter-linux-arm64 .
	@echo "Linux build complete!"

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p dist
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -o dist/gin-starter-windows-amd64.exe .
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build -o dist/gin-starter-windows-arm64.exe .
	@echo "Windows build complete!"

# Build for macOS (Darwin)
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p dist
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -o dist/gin-starter-darwin-amd64 .
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o dist/gin-starter-darwin-arm64 .
	@echo "macOS build complete!"

# Build for specific platform (usage: make build-platform OS=<os> ARCH=<arch>)
build-platform:
	@echo "Building for $(OS)/$(ARCH)..."
	@mkdir -p dist
	GOOS=$(OS) GOARCH=$(ARCH) go build -o dist/gin-starter-$(OS)-$(ARCH)$(if $(filter windows,$(OS)),.exe,) .
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf dist
	@echo "Clean complete!"
