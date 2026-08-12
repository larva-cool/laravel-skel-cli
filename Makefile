BINARY  := laravel-skel-cli
VERSION ?= v1.0.0
LDFLAGS := -s -w -X laravel-skel-cli/cmd.version=$(VERSION)

# 当前平台
.PHONY: build
build:
	mkdir -p dist
	go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY) .

# 跨平台打包（macOS/Windows/Linux，arm64+amd64）
.PHONY: release
release:
	mkdir -p dist
	GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe .

.PHONY: clean
clean:
	rm -rf dist

.PHONY: test
test:
	go vet ./...
	go build ./...
