PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := internal/pb

.PHONY: generate build-validator build-gateway build test docker-build docker-run docker-stop

generate:
	mkdir -p $(GO_OUT)
	protoc \
		-I$(PROTO_DIR) \
		-Iproto/googleapis \
		--go_out=paths=source_relative:$(GO_OUT) \
		--go-grpc_out=paths=source_relative:$(GO_OUT) \
		--grpc-gateway_out=paths=source_relative:$(GO_OUT) \
		$(PROTO_FILES)

build-validator:
	docker build -f Dockerfile.email-validator -t email-validator:latest .

build-gateway:
	docker build -f Dockerfile.gateway -t email-gateway:latest .

build: build-validator build-gateway

test:
	go test ./...

docker-build:
	make generate
	docker compose build

docker-run:
	docker compose up --build

docker-stop:
	docker compose down
