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

build:
	docker build --tag ngoctd/c72-server .   

push:
	docker push docker.io/ngoctd/c72-server:latest
