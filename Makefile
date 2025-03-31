gomodule:=github.com/ngoctd314/c72-api-server

dev:
	go run ./cmd/api

migrate:
	go run ./cmd/migrate $(filter-out $@,$(MAKECMDGOALS))

migrate-create:
	migrate create -ext sql -dir migrations -seq ${name}

build:
	docker build --tag ngoctd/c72-server .   

push:
	docker push docker.io/ngoctd/c72-server:latest

testfunc:
	go test -v -count=1 -run ${func} ${gomodule}/${pkg}

%:
	@:
