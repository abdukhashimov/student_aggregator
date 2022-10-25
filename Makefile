.PHONY: run
run:
	go run cmd/student_aggregator/main.go http -port=8484

.PHONY: mongo_dev_start
mongo_dev_start:
	docker compose -f docker-compose.dev.yml up -d

.PHONY: mongo_dev_stop
mongo_dev_stop:
	docker compose -f docker-compose.dev.yml down

.PHONY: mongo_dev_remove_db
mongo_dev_remove_db:
	docker compose -f docker-compose.dev.yml down --volumes
