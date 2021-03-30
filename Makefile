test:
	go test ./... -v

build: get
	go build -ldflags "-X 'github.com/xos/rssbot/config.commit=`git rev-parse --short HEAD`' -X 'github.com/xos/rssbot/config.date=`date`'"

get:
	go mod download

run:
	go run .