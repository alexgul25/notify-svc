run-svc:
	@echo "🚀  Применение миграций и запуск Notify Service..."
	@go run ./cmd/migrator/main.go && go run ./cmd/svc-starter/main.go; exit 0