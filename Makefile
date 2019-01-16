test:
	# The item just for travis-ci :)

build: get
	go build .

get:
	go mod download