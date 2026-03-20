tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

validate ::
	@echo "🔍 Running validation checks..."
	@echo ""
	@echo "1️⃣  Checking code formatting..."
	@FILES=$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*' 2>/dev/null) 2>/dev/null); \
	if [ -n "$$FILES" ]; then \
		echo "❌ Files need formatting:"; \
		echo "$$FILES"; \
		echo ""; \
		echo "Run: gofmt -w ."; \
		exit 1; \
	fi
	@echo "✅ Code formatting is correct"
	@echo ""
	@echo "2️⃣  Running go vet (static analysis)..."
	@go vet ./... 2>&1 || (echo "❌ go vet found issues" && exit 1)
	@echo "✅ go vet passed"
	@echo ""
	@echo "3️⃣  Checking for common mistakes with staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./... 2>&1 || (echo "❌ staticcheck found issues" && exit 1); \
		echo "✅ staticcheck passed"; \
	else \
		echo "⚠️  staticcheck not installed (optional - install: go install honnef.co/go/tools/cmd/staticcheck@latest)"; \
	fi
	@echo ""
	@echo "4️⃣  Running tests with race detector and coverage..."
	@go test -race -count=1 ./... -coverprofile=coverage.out -covermode=atomic 2>&1 || (echo "❌ Tests failed" && exit 1)
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "✅ All tests passed (coverage: $${COVERAGE}%)"
	@echo ""
	@echo "5️⃣  Verifying dependencies..."
	@go mod verify >/dev/null 2>&1 || (echo "❌ Dependency verification failed" && exit 1)
	@echo "✅ Dependencies verified"
	@echo ""
	@echo "6️⃣  Checking for dependency issues..."
	@cp go.mod go.mod.bak && cp go.sum go.sum.bak
	@go mod tidy >/dev/null 2>&1
	@if ! diff -q go.mod go.mod.bak >/dev/null 2>&1 || ! diff -q go.sum go.sum.bak >/dev/null 2>&1; then \
		mv go.mod.bak go.mod && mv go.sum.bak go.sum; \
		echo "⚠️  Dependencies need tidying (run: make tidy)"; \
	else \
		rm go.mod.bak go.sum.bak; \
		echo "✅ Dependencies are tidy"; \
	fi
	@echo ""
	@echo "7️⃣  Checking for security vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./... 2>&1 || (echo "❌ Security vulnerabilities found" && exit 1); \
		echo "✅ No known vulnerabilities"; \
	else \
		echo "⚠️  govulncheck not installed (optional - install: go install golang.org/x/vuln/cmd/govulncheck@latest)"; \
	fi
	@echo ""
	@echo "✨ All validation checks passed!"

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down
