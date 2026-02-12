# ===== config =====
ANALYTICS_PATH=/Users/Alex/IdeaProjects/url-shortener/services/analytics-service/cmd/main.go
CACHE_PATH=/Users/Alex/IdeaProjects/url-shortener/services/cache-service/cmd/main.go
GATEWAY_PATH=/Users/Alex/IdeaProjects/url-shortener/services/api-gateway/cmd/main.go
URL_PATH=/Users/Alex/IdeaProjects/url-shortener/services/url-service/cmd/main.go
# ==================

#export GOOSE_DRIVER=$(DB_DRIVER)
#export GOOSE_DBSTRING=$(DB_DSN)

run-analytics:
	go run $(ANALYTICS_PATH)
run-cache:
	go run $(CACHE_PATH)
run-gateway:
	go run $(GATEWAY_PATH)
run-url:
	go run $(URL_PATH)
run-compose:
	docker compose -f docker-compose.dev.yml up -d
