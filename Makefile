EXPORT := export $$(cat .env|xargs)
APP_NAME := karmabot

.PHONY: test
test:
	APP_NAME=$(APP_NAME) go test ./...

.PHONY: run
run:
	$(EXPORT) && go run ./cmd/$(APP_NAME)

.PHONY: lint
lint:
	golangci-lint run
