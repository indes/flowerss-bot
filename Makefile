app_name = flowerss-bot

VERSION=$(shell git describe --tags --always)
DATA=$(shell date)
COMMIT=$(shell git rev-parse --short HEAD)
test:
	go test ./... -v

all: build

build: get
	go build -ldflags \
		"-X 'github.com/indes/flowerss-bot/config.commit=$(COMMIT)' \
		-X 'github.com/indes/flowerss-bot/config.date=$(DATA)' \
		-X 'github.com/indes/flowerss-bot/config.version=$(VERSION)'" -o $(app_name)

get:
	go mod download

run:
	go run .

clean:
	rm flowerss-bot