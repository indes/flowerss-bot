test:
	go test ./... -v

build: get
	go build -ldflags "-X 'github.com/yangon99/flowerss-bot/config.commit=`git rev-parse --short HEAD`' -X 'github.com/yangon99/flowerss-bot/config.date=`date`'"

get:
	go mod download

run:
	go run .