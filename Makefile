MIGRATE_RUN=go run db/scripts/migrate.go

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=create_table_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	$(MIGRATE_RUN) up

migrate-down:
	$(MIGRATE_RUN) down 1

migrate-version:
	$(MIGRATE_RUN) version

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Usage: make migrate-force version=<version>"; \
		exit 1; \
	fi
	$(MIGRATE_RUN) force $(version)

reset-db:
	@echo "⚠️  WARNING: Resetting DB to version 0 (dev only)..."
	go run scripts/migrate.go force 1
	go run scripts/migrate.go down
	go run scripts/migrate.go up
	@echo "✅ Database reset and migrated to latest"