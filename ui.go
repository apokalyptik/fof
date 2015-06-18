package main

//go:generate sass $PWD/www/css/style.scss $PWD/www/css/style.css
//go:generate browserify -t reactify -t uglifyify $PWD/www/js/app/main.jsx -o $PWD/www/js/production.js
//go:generate $GOPATH/bin/go-bindata -ignore=$PWD/www/css/.*.scss -prefix $PWD/ $PWD/www/ $PWD/www/js $PWD/www/css
