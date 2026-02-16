SSO_PG_DSN=postgres://myuser:mypassword@localhost:5555/postgres?sslmode=disable
SSO_MIGRATION=sql/goose/sso/migrations/

.PHONY: generate-private test-integrate sso-migrate-create sso-migrate-up sso-migrate-down sso-migrate-status sso-migrate-reset run

# Default target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  docker-up                            Start all containers"
	@echo "  docker-down                          Stop all containers"
	@echo "  migrate-create name=<table_name>     Create new migration file"
	@echo "  migrate-up                           Apply all pending migrations"
	@echo "  migrate-down                         Rollback last migration"
	@echo "  migrate-status                       Show migrations status"
	@echo "  migrate-reset                        Rollback all migrations"
	@echo "  test                                 Start tests"
	@echo "  test-integrate                       Start intergration tests"

generate-private:
	go run scripts/gen-private.go

run:
	go run cmd/sso/main.go -config configs/sso.yaml

# Create new migration file: make sso-migrate-create name=name_table
sso-migrate-create:
	goose -dir $(SSO_MIGRATION) create $(name) sql

sso-migrate-up:
	goose -dir $(SSO_MIGRATION) postgres "$(SSO_PG_DSN)" up

sso-migrate-down:
	goose -dir $(SSO_MIGRATION) postgres "$(SSO_PG_DSN)" down

sso-migrate-status:
	goose -dir $(SSO_MIGRATION) postgres "$(SSO_PG_DSN)" status

sso-migrate-reset:
	goose -dir $(SSO_MIGRATION) postgres "$(SSO_PG_DSN)" reset

# Start tests
test:
	go test -v -count=1 ./...

test-integrate:
	go test -tags=integration -v ./...
