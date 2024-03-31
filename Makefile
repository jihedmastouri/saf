.PHONY: init
init:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go mod tidy

.PHONY: sqlgen
sqlgen:
	@sqlc generate -f ./postgres/sqlc.yaml
