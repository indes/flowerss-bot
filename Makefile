test:
	# The item just for travis-ci :)

build: get
	go build -ldflags "-X 'github.com/indes/flowerss-bot/config.commit=`git rev-parse --short HEAD`' -X 'github.com/indes/flowerss-bot/config.date=`date`'"

get:
	go mod download
