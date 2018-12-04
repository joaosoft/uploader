env:
	docker-compose up -d uploader.postgres uploader.redis uploader.rabbitmq
run:
	go run ./main.go

build:
	go build .

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*