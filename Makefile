dsn:="mysql://root:secret@tcp(localhost:3306)/tag_scan"

migrate-create:
	migrate create -ext sql -dir migrations -seq ${name}

migrate-up:
	migrate -path ./migrations/ -database ${dsn} up

migrate-down:
	migrate -path ./migrations/ -database ${dsn} down

migrate-force:
	migrate -path ./migrations/ -database ${dsn} force ${version}

dev:
	go run ./cmd

