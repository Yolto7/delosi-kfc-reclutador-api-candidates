APP_NAME := api-candidates

# Lambda paths
CMD_LAMBDAS_MAIN := cmd/lambdas/main.go

# Directory where binaries and ZIPs are placed
BUILD_DIR := bin

.PHONY: clean deps build deploy-dev deploy-prod remove-dev remove-prod

clean:
	@echo "🧹 Cleaning binaries and generated files..."
	@go clean
	@rm -rf ./$(BUILD_DIR)
	@rm -rf ./vendor
	@echo "✅ Clean complete"

deps:
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed"

build-main:
	@echo "🚀 Building main lambdas..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap $(CMD_LAMBDAS_MAIN)
	@zip -j $(BUILD_DIR)/main.zip bootstrap
	@rm -f bootstrap
	@echo "✅ Main lambda built and zipped"

build: deps build-main
	@echo "🔨 Building $(APP_NAME) complete..."

deploy-dev: build
	@echo "☁️ Deploying $(APP_NAME) to dev stage..."
	sls deploy --stage dev
	@echo "✅ Deployment to dev complete"

deploy-prod: build
	@echo "☁️ Deploying $(APP_NAME) to prod stage..."
	sls deploy --stage prod
	@echo "✅ Deployment to prod complete"

remove-dev:
	@echo "🚫 Removing $(APP_NAME) from dev stage..."
	sls remove --stage dev
	@echo "✅ Removal from dev complete"

remove-prod:
	@echo "🚫 Removing $(APP_NAME) from prod stage..."
	sls remove --stage prod
	@echo "✅ Removal from prod complete"
