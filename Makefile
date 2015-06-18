all:
	npm install
	go generate
	go build
dev:
	sass www/css/style.scss www/css/style.css
	browserify -t reactify --debug www/js/app.jsx -o www/js/development.js
