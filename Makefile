SSO_PG_DSN=postgres://myuser:mypassword@localhost:5555/postgres?sslmode=disable
SSO_MIGRATION=sql/goose/sso/migrations/

.PHONY: generate-private test-integrate

generate-private:
	go run scripts/gen-private.go

test-integrate:
	go test -tags=integration -v ./...

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
