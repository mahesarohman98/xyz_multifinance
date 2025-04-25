include .env
export

.PHONY: mycli
mysql:
	docker exec -it mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)
	
test: test-up wait-db
	@./scripts/test.sh .env
	@$(MAKE) test-down

test-up:
	docker compose -f docker-compose.test.yml up -d --build

test-down:
	docker compose -f docker-compose.test.yml down -v

wait-db:
	@until docker exec mysql mysqladmin ping -h "127.0.0.1" --silent; do \
		sleep 1; \
	done
