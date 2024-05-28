download:
	@echo Download go.mod dependencies
	@go mod download

install/tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
	@npm i

run/live:
	@air

create/db:
	@touch db/dev.db
	migrate -database sqlite3://db/dev.db -path db/migrations up

create/migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

push/migrations:
	migrate -database sqlite3://db/dev.db -path db/migrations up
