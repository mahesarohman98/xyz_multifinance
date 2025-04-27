include .env
export

.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh customer src/internal/customer/ports ports
	@./scripts/openapi-http.sh creditlimit src/internal/creditlimit/ports ports

.PHONY: mycli
mysql:
	docker exec -it mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)

dev-up:
	docker compose up -d --build

dev-down:
	docker compose down -v
	
test: test-up wait-db
	@./scripts/test.sh .test.env
	@$(MAKE) test-down

test-up:
	docker compose -f docker-compose.test.yml up -d --build

test-down:
	docker compose -f docker-compose.test.yml down -v

wait-db:
	@until docker exec mysql mysqladmin ping -h "127.0.0.1" --silent; do \
		sleep 1; \
	done
