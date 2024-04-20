.DEFAULT_GOAL := help

.PHONY: run
run: ## 実行する
	go run .

.PHONY: memプロファイルを見る
memp: ## 実行する
	go tool pprof mem.pprof

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
