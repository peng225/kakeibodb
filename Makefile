KAKEIBODB=./kakeibodb
PASSWORD?=
GO_FILES:=$(shell find . -type f -name '*.go' -print)
COMMON_OPTIONS:=--dbname testdb -u test
COMMON_OPTIONS_WO_USER:=--dbname testdb

$(KAKEIBODB): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: test
test: $(KAKEIBODB)
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(KAKEIBODB)

.PHONY: e2e-test
e2e-test: $(KAKEIBODB)
# Load events.
	$(KAKEIBODB) event load -d test/event $(COMMON_OPTIONS)
	CREDIT_EVENT_ID=$$($(KAKEIBODB) event list  $(COMMON_OPTIONS) | grep "クレジット" | awk '{print $$1}'); \
	echo $${CREDIT_EVENT_ID}; \
	$(KAKEIBODB) event load --credit --parentEventID $${CREDIT_EVENT_ID} -f test/credit/cmeisai1.csv $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
# Tag create, add, remove, and delete.
	$(KAKEIBODB) tag create -n foo $(COMMON_OPTIONS)
	$(KAKEIBODB) tag list $(COMMON_OPTIONS)
	$(KAKEIBODB) event addTag -e 1 -t foo $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) event removeTag -e 1 -t foo $(COMMON_OPTIONS)
	$(KAKEIBODB) tag delete -t 1 $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) tag list $(COMMON_OPTIONS)
# Set user name by env.
	KAKEIBODB_USER=test $(KAKEIBODB) event list $(COMMON_OPTIONS_WO_USER)

.PHONY: test-setup
test-setup:
	mysql -h 127.0.0.1 --port 3306 -B -u root -p$(PASSWORD) < test/setup.sql

.PHONY: test-clean
test-clean:
	mysql -B -u root -p $(PASSWORD) < test/clean.sql
