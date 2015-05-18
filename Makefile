all:
	npm install
	go generate
	go build
