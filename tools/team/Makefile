all:
	npm install
	go generate
	go build
dev:
	watchify -v \
		-t reactify \
		--debug www/js/app/main.jsx \
		-o www/js/development.js
