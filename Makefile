gomodule:=github.com/ngoctd314/c72-api-server

dev:
	go run ./cmd/api
	# docker run -p 5080:5080 ngoctd/c72-server

migrate:
	go run ./cmd/migrate $(filter-out $@,$(MAKECMDGOALS))

migrate-create:
	migrate create -ext sql -dir migrations -seq ${name}

build:
	docker build --tag ngoctd/c72-server .   

push:
	docker push docker.io/ngoctd/c72-server:latest

run:
	docker run -dp 5080:5080 ngoctd/c72-server

testfunc:
	go test -v -count=1 -run ${func} ${gomodule}/${pkg}

%:
	@:
