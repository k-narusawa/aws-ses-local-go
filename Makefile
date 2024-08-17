dev:
	docker-compose up

up: 
	docker-compose up -d

down:
	docker-compose down

test-ci:
	go test -timeout=3s -race -count=10 -failfast -shuffle=on -short ./... -coverprofile=./cover.short.profile -covermode=atomic -coverpkg=./...
	go test -timeout=10s -race -count=1 -failfast  -shuffle=on ./... -coverprofile=./cover.long.profile -covermode=atomic -coverpkg=./...

test: 
	go test -v ./... -count=1 -cover -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out -o out/coverage.html

build:
	@ printf "Building aplication... "
	@ go build \
		-trimpath  \
		-o engine \
		./app/
	@ echo "done"

dev-build:
	docker build -f Dockerfile.dev -t go-app-template .
