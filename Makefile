KAKEIBODB := ./kakeibodb
GO_FILES := $(shell find . -type f -name '*.go' -print)

$(KAKEIBODB): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: generate
generate:
	sqlc generate -f internal/repository/mysql/sqlc/sqlc.yaml

.PHONY: test
test: $(KAKEIBODB)
	go test -v ./...

.PHONYE: test-setup
test-setup:
#	docker run -d --rm --env MYSQL_ALLOW_EMPTY_PASSWORD=yes -p 3307:3306 mysql:8
	mysql -h 127.0.0.1 --port 3307 -B -u root < internal/test/setup.sql

.PHONYE: test-cleanup
test-cleanup:
	mysql -h 127.0.0.1 --port 3307 -B -u root < internal/test/cleanup.sql


.PHONY: clean
clean:
	rm -f $(KAKEIBODB)
