# NijiChat Backend Makefile

.PHONY: help
help: ## 显示帮助信息
	@echo "NijiChat Backend 开发命令"
	@echo "=========================="
	@echo "  build    构建WASM应用"
	@echo "  dev      启动开发服务器"
	@echo "  deploy   部署应用"
	@echo "  clean    清理构建文件"

.PHONY: build
build: ## 构建WASM应用
	@npm run build

.PHONY: dev
dev: ## 启动开发服务器
	@npm run dev

.PHONY: deploy
deploy: ## 部署应用
	@npm run deploy

.PHONY: clean
clean: ## 清理构建文件
	@rm -rf build/ && go clean