.PHONY: local build server check fmt clean gittag tag

APP=sleep_admin

ENV = $(firstword $(filter-out $@,$(MAKECMDGOALS)))

BIN_DIR := ./bin
IMAGE ?= app-client

Linux := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_FLAGS := -tags="sonic" -ldflags "-s -w"
## 构建mac版本,需要安装upx
# 多个 tag 的正确写法: 多个 tag 之间用逗号或空格分隔
#go build -tags=sonic,debug ./...
#go build -tags="sonic debug" ./...


local:
	go run -tags="sonic" . server -c dev

## 构建linux版本,需要安装upx
build:
	@mkdir -p $(BIN_DIR)
	$(Linux) go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP) . && upx -9 $(BIN_DIR)/$(APP)

server:
	$(BIN_DIR)/${APP} server  -c $(ENV)


LINT_TARGETS ?= ./...
tool: ## Lint Go code with the installed golangci-lint
	@ echo "▶️ golangci-lint run"
	golangci-lint run $(LINT_TARGETS)
	@ echo "✅ golangci-lint run"

## govulncheck 检查漏洞 `go install golang.org/x/vuln/cmd/govulncheck@latest`
## gosec 检查安全漏洞 `go install github.com/securego/gosec/v2/cmd/gosec@latest`
check:
	govulncheck ./...
	gosec ./...
## betteralign 优化结构体字段排序和内存布局 `go install github.com/dkorunic/betteralign/cmd/betteralign@latest`
better:
	betteralign ./...
	betteralign -apply ./...
## gofumpt 格式化代码 `go install mvdan.cc/gofumpt@latest`
fmt:
	gofumpt -l -w .
	go mod tidy

clean:
	@echo "Cleaning..."
	@rm -rf ./${APP}

gittag:
	git tag --sort=-version:refname | head -1

tag:
	@current=$$(grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' version.go | head -n1 | tr -d 'v'); \
	if [ -z "$$current" ]; then echo "version not found in version.go"; exit 1; fi; \
	maj=$$(echo $$current | cut -d. -f1); \
	min=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	newpatch=$$(expr $$patch + 1); \
	new="v$$maj.$$min.$$newpatch"; \
	printf "Bump: v%s -> %s\n" "$$current" "$$new"; \
	sed -E -i.bak 's/(const Version = ")([^"]+)(")/\1'"$$new"'\3/' version.go; \
	git add version.go; \
	git commit -m "chore(release): $$new"; \
	printf "Release: %s\n" "$$new"; \
	git push ydsd HEAD; \
	git tag -a "$$new" -m "release $$new"; \
	printf "Tag: %s\n" "$$new"; \
	git push ydsd "$$new"; \
	printf "Done\n"
	rm -f version.go.bak