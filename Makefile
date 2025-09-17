GOFLAGS += -trimpath

tmpdir = tmp

.PHONY: test
test:
	@mkdir -p "$(tmpdir)/reports"
	@echo "testing..."
	@go test $(GOFLAGS) -coverprofile "$(tmpdir)/reports/coverage.out" ./...
	@echo "generating coverage reports..."
	@go tool cover -html "$(tmpdir)/reports/coverage.out" -o "$(tmpdir)/reports/coverage.html"
	@echo "done"

.PHONY: clean
clean:
	@rm -rf "$(tmpdir)"
