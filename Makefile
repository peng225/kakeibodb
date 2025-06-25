KAKEIBODB := ./kakeibodb
PASSWORD ?=
GO_FILES := $(shell find . -type f -name '*.go' -print)
COMMON_OPTIONS := --dbname testdb -u test
COMMON_OPTIONS_WO_USER := --dbname testdb

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
	$(KAKEIBODB) event load -d internal/test/event $(COMMON_OPTIONS)
	CREDIT_EVENT_ID=$$($(KAKEIBODB) event list  $(COMMON_OPTIONS) | grep "クレジット" | awk '{print $$1}'); \
	echo $${CREDIT_EVENT_ID}; \
	$(KAKEIBODB) event load --credit --parentEventID $${CREDIT_EVENT_ID} -f internal/test/credit/cmeisai1.csv $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
# Tag create, add, remove, and delete.
	$(KAKEIBODB) tag create -t foo $(COMMON_OPTIONS)
	$(KAKEIBODB) tag create -t bar $(COMMON_OPTIONS)
	$(KAKEIBODB) tag list $(COMMON_OPTIONS)
	$(KAKEIBODB) event addTag --eventID 1 --tagNames foo,bar $(COMMON_OPTIONS)
# Idempotency check for addTag.
	$(KAKEIBODB) event addTag --eventID 1 --tagNames foo $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) event removeTag --eventID 1 -t foo $(COMMON_OPTIONS)
	$(KAKEIBODB) event removeTag --eventID 1 -t bar $(COMMON_OPTIONS)
	$(KAKEIBODB) tag delete --tagID 1 $(COMMON_OPTIONS)
	$(KAKEIBODB) tag delete --tagID 2 $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) tag list $(COMMON_OPTIONS)
# Set user name by env.
	KAKEIBODB_USER=test $(KAKEIBODB) event list $(COMMON_OPTIONS_WO_USER)
# Pattern lifecycle.
	$(KAKEIBODB) pattern list $(COMMON_OPTIONS)
	$(KAKEIBODB) tag create -t fruit $(COMMON_OPTIONS)
	$(KAKEIBODB) tag create -t yellow $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern create -k "バナ" $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern list $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern addTag --patternID 1 --tagNames fruit,yellow $(COMMON_OPTIONS)
# Idempotency check for addTag.
	$(KAKEIBODB) pattern addTag --patternID 1 --tagNames fruit,yellow $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern list $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) event applyPattern --from "2022-01-04" --to "2022-02-03" $(COMMON_OPTIONS)
# Idempotency check for applyPattern.
	$(KAKEIBODB) event applyPattern --from "2022-01-04" --to "2022-02-03" $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern removeTag --patternID 1 -t fruit $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern list $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern delete --patternID 1 $(COMMON_OPTIONS)
	$(KAKEIBODB) pattern list $(COMMON_OPTIONS)
# Split test
	$(KAKEIBODB) tag create -t candy $(COMMON_OPTIONS)
	$(KAKEIBODB) event addTag --eventID 10 --tagNames "candy" $(COMMON_OPTIONS)
	$(KAKEIBODB) event split --eventID 10 --date 2021-12-04 --money -30 --desc "はちみつのど飴" $(COMMON_OPTIONS)
	$(KAKEIBODB) event split --eventID 10 --date 2021/12/05 --money -30 --desc "きんかんのど飴" $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)
# Split with auto eventID detection
	KAKEIBODB_SPLIT_BASE_TAG_NAME="candy" $(KAKEIBODB) event split --date 2021-12-06 --money -40 --desc "ミルク飴" $(COMMON_OPTIONS)
	$(KAKEIBODB) event list $(COMMON_OPTIONS)

.PHONY: test-setup
test-setup:
	mysql -h 127.0.0.1 --port 3306 -B -u root -p$(PASSWORD) < internal/test/setup.sql

.PHONY: test-clean
test-clean:
	mysql -B -u root -p$(PASSWORD) < internal/test/clean.sql
