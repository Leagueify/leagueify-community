build:
	docker build -t leagueify .

build-dev:
	docker build --target dev -t leagueify-dev .

clean: stop
	docker compose -f docker-compose.prod.yml down -v

clean-dev: stop-dev
	docker compose -f docker-compose.dev.yml down -v

format:
	go fmt ./...

init:
	go get ./...

prep: format tidy vet

start: build
	docker compose -f docker-compose.prod.yml up

start-detached: build
	docker compose -f docker-compose.prod.yml up -d

start-dev: build-dev
	docker compose -f docker-compose.dev.yml up

start-dev-detached: build-dev
	docker compose -f docker-compose.dev.yml up -d

stop:
	docker compose -f docker-compose.prod.yml down

stop-dev:
	docker compose -f docker-compose.dev.yml down

test:
	mkdir -p testCoverage
	go test ./... -cover -coverprofile=testCoverage/report.out
	go tool cover -html=testCoverage/report.out -o testCoverage/report.html

tidy:
	go mod tidy

vet:
	go vet ./...
