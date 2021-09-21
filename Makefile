app_name = flowerss-bot

test:
	go test ./... -v

build: get
	go build -ldflags "-X 'github.com/indes/flowerss-bot/config.commit=`git rev-parse --short HEAD`' -X 'github.com/indes/flowerss-bot/config.date=`date`'" -o $(app_name)

get:
	go mod download

run:
	go run .

clean:
	rm flowerss-bot