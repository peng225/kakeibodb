BUILD_TARGET=kakeibodb

GO_FILES:=$(shell find . -type f -name '*.go' -print)
MINIO_DATAPATH:=~/minio/data

$(BUILD_TARGET): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: test
test: $(BUILD_TARGET)
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(BUILD_TARGET)