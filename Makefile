KAKEIBODB := ./kakeibodb
GO_FILES := $(shell find . -type f -name '*.go' -print)

$(KAKEIBODB): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: generate
generate:
	sqlc generate -f internal/repository/mysql/query/sqlc.yaml

.PHONY: test
test: $(KAKEIBODB)
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(KAKEIBODB)
