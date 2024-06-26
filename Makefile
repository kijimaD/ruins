.DEFAULT_GOAL := help

.PHONY: run
run: ## 実行する。スクショのキーを指定している
	EBITENGINE_SCREENSHOT_KEY=1 go run .

.PHONY: test
test: ## テストを実行する
	go test ./... -v

.PHONY: build
build: ## ビルドする
	./scripts/build.sh

.PHONY: vrt
vrt: ## 各ステートでスクショを取得する
	./scripts/vrt.sh

.PHONY: memp
memp: ## 実行毎に保存しているプロファイルを見る
	go tool pprof mem.pprof

.PHONY: pprof
pprof: ## サーバ経由で取得したプロファイルを見る。起動中でなければならない
	go build .
	go tool pprof ruins "http://localhost:6060/debug/pprof/profile?seconds=10"

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
