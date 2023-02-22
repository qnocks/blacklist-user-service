.SILENT:

build:
	go build -o ./.bin/app cmd/app/main.go

run: build
	./.bin/app

test:
	go test ./cmd/... ./internal/... -race -coverprofile=cover.out ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out | grep "total"

lint:
	golangci-lint run

swagger:
	swag init -g internal/app/app.go

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock.go

docker:
	docker build -t blacklist-app .

compose:
	docker compose up
