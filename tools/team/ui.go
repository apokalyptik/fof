package main

//go:generate env PRODUCTION=1 webpack -p
//go:generate /bin/bash ./cachebuster.sh
//go:generate $GOPATH/bin/go-bindata -ignore=$PWD/www/css/.*.scss -prefix $PWD/ $PWD/www/ $PWD/www/js $PWD/www/css $PWD/www/fonts
//go:generate git checkout -- www/index.html
