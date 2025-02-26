GIT_VERSION := $(shell git describe --tags --always --dirty)
BINARY_NAME := ionlog-viewer-${GIT_VERSION}
MAIN_PATH := ./cmd/main.go

.PHONY: build/develop
build/develop:
	@echo "Building development verison"
	go build -o bin/${BINARY_NAME} ${MAIN_PATH}


.PHONY: build/production
build/production:
	@echo "Building production version"
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}.exe ${MAIN_PATH}
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME} ${MAIN_PATH}


.PHONY: clean
clean:
	@echo "Cleaning bin/ folder"
	@rm -rf bin/


.PHONY: audit
audit:
	go fmt ./...
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go test -vet=off ./...

