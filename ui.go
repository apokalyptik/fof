package main

//go:generate sass $PWD/www/style.scss $PWD/www/style.css
//go:generate $GOPATH/bin/go-bindata -ignore=$PWD/www/.*.scss -prefix $PWD/ $PWD/www/
