default: build

clean:
	rm -f cmd/server/server
	rm -f cmd/todo/todo

build:
	cd cmd/server; \
	go build
	cd cmd/todo; \
	go build

deps:
	cd cmd/server; \
	go get

migrate:
	./cmd/server/server --config config.yaml migratedb

test:
	./cmd/server/server --config config.yaml server & \
	pid=$$!; \
	go test; \
	kill $$pid
