default: build

deps:
	go get

build:
	cd cmd/server; \
	go build
	cd cmd/todo; \
	go build


test:
	./cmd/server/server --config config.yaml server & \
	pid=$$!; \
	go test; \
	kill $$pid
