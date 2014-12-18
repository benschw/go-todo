default: build

build:
	cd cmd/server; \
	go build
	cd cmd/todo; \
	go build

deps:
	go get

migrate:
	./cmd/server/server --config config.yaml migratedb

test:
	./cmd/server/server --config config.yaml server & \
	pid=$$!; \
	go test; \
	kill $$pid
