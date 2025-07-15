.DEFAULT_GOAL := help

.PHONY: run
run: ## 実行する。スクショのキーを指定している
	EBITENGINE_SCREENSHOT_KEY=1 go run .

.PHONY: test
test: ## テストを実行する
	go test -v -cover -race ./...

.PHONY: build
build: ## ビルドする
	./scripts/build.sh

.PHONY: vrt
vrt: ## 各ステートでスクショを取得する
	./scripts/vrt.sh

.PHONY: fmt
fmt: ## フォーマットする
	goimports -w .

.PHONY: lint
lint: ## Linterを実行する
	@echo "Checking formatting..."
	@gofmt -l ./... | grep -E '\.go$$' && echo "Code needs formatting. Run 'make fmt'" && exit 1 || echo "Code is properly formatted"
	@echo "Checking imports..."
	@goimports -l ./... | grep -E '\.go$$' && echo "Imports need formatting. Run 'make goimports'" && exit 1 || echo "Imports are properly formatted"
	@echo "Running golangci-lint..."
	@golangci-lint run -v ./...

.PHONY: tools-install
tools-install: ## 開発ツールをインストールする
	@go install golang.org/x/tools/cmd/goimports@latest
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)

.PHONY: check
check: test build fmt lint ## 一気にチェックする

# ================

.PHONY: memp
memp: ## 実行毎に保存しているプロファイルを見る
	go tool pprof mem.pprof

.PHONY: pprof
pprof: ## サーバ経由で取得したプロファイルを見る。起動中でなければならない
	go build .
	go tool pprof ruins "http://localhost:6060/debug/pprof/profile?seconds=10"

# ================

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
