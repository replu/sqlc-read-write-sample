
.PHONY: migrate-dryrun
migrate-dryrun:
	docker compose run --rm migrator \
		bash -c "mysqldef -u root -p password -P 3306 -h mysql-primary sample --dry-run < ./db/schema.sql"

.PHONY: migrate
migrate:
	docker compose run --rm migrator \
		bash -c "mysqldef -u root -p password -P 3306 -h mysql-primary sample < ./db/schema.sql"
