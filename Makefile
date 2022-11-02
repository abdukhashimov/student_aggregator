.PHONY: run
run:
	go run cmd/main.go http --port=8484

.PHONY: dev_environment_start
dev_environment_start:
	docker compose -f docker-compose.dev.yml up -d

.PHONY: dev_environment_stop
dev_environment_stop:
	docker compose -f docker-compose.dev.yml down

.PHONY: dev_environment_remove
dev_environment_remove:
	docker compose -f docker-compose.dev.yml down --volumes

.PHONY: swagger
swagger:
	swag init -g internal/transport/handlers/server.go
