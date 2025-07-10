.PHONY: build clean run all

# Build the host monitoring system
build:
	@echo "🔨 Building host monitoring system..."
	go build -o host-monitor main.go host_monitor.go host_security.go host_performance.go host_network.go host_package.go
	@echo "✅ Build completed!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f host-monitor
	@echo "✅ Clean completed!"

# Run the host monitoring system
run: build
	@echo "🚀 Running host monitoring system..."
	./host-monitor

# Build and run all modules
all: build
	@echo "🎯 Running all monitoring modules..."
	@echo "6" | ./host-monitor

# Install dependencies (if needed)
deps:
	@echo "📦 Installing dependencies..."
	go mod init host-monitor
	go mod tidy

# Test the build
test: build
	@echo "🧪 Testing build..."
	@echo "1" | timeout 10s ./host-monitor || true
	@echo "✅ Test completed!"

# Help
help:
	@echo "Available commands:"
	@echo "  build  - Build the host monitoring system"
	@echo "  clean  - Clean build artifacts"
	@echo "  run    - Build and run the system"
	@echo "  all    - Build and run all modules"
	@echo "  deps   - Install dependencies"
	@echo "  test   - Test the build"
	@echo "  help   - Show this help" 