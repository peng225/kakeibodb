KAKEIBODB := ./kakeibodb
GO_FILES := $(shell find . -type f -name '*.go' -print)

BINDIR := bin
GOLANGCI_LINT_VERSION := v2.2.2
GOLANGCI_LINT := $(BINDIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)

$(KAKEIBODB): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

$(BINDIR):
	mkdir -p $@

.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

$(GOLANGCI_LINT): | $(BINDIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b . $(GOLANGCI_LINT_VERSION)
	mv golangci-lint $(GOLANGCI_LINT)

.PHONY: generate
generate:
	sqlc generate -f internal/repository/mysql/sqlc/sqlc.yaml

.PHONY: lint
lint: | $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: test
test: $(KAKEIBODB)
	go test -count=1 -v ./internal/...

.PHONY: test-setup
test-setup:
#	docker run -d --rm --env MYSQL_ALLOW_EMPTY_PASSWORD=yes -p 3307:3306 mysql:8
	mysql -h 127.0.0.1 --port 3307 -B -u root < internal/test/setup.sql

.PHONY: test-cleanup
test-cleanup:
	mysql -h 127.0.0.1 --port 3307 -B -u root < internal/test/cleanup.sql

.PHONY: clean
clean:
	rm -f $(KAKEIBODB)
