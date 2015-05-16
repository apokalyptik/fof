package main

//go:generate sass $PWD/www/css/style.scss $PWD/www/css/style.css
//go:generate $GOPATH/bin/go-bindata -ignore=$PWD/www/css/.*.scss -prefix $PWD/ $PWD/www/ $PWD/www/js $PWD/www/css
