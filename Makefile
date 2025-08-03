.PHONY: server
server: site db-local
	go run main.go start

.PHONY: worker db-local
worker:
	go run main.go start --worker

.PHONY: db-local
db-local:
	docker compose -f docker-compose.dev.yaml up -d

.PHONY: site
site:
	cd site && npm run build

