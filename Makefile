.PHONY: deps build test lint run export docker-docs docker-cli clean

APP := gitpulse
BIN := bin/$(APP)

deps:
	go mod download
	go mod tidy

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o $(BIN) .

test:
	go test ./... -race -count=1

lint:
	go vet ./...

run: build
	./$(BIN)

export: build
	./$(BIN) export -f html -o sample
	./$(BIN) export -f json -o sample
	./$(BIN) export -f csv -o sample
	./$(BIN) export -f pdf -o sample

docker-docs:
	docker compose up --build -d gitpulse-docs

docker-cli:
	docker build --target runtime -t gitpulse:cli .

clean:
	rm -rf bin/ reports/gitpulse-* coverage.out
