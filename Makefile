KAKEIBODB := ./kakeibodb
GO_FILES := $(shell find . -type f -name '*.go' -print)

$(KAKEIBODB): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: test
test: $(KAKEIBODB)
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(KAKEIBODB)
