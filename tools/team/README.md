# raidbot
Federation of Fathers Slack Raid integration code

## Getting Team Tool running
* Install Go
 * https://golang.org/doc/install
 * Make sure to set $GOROOT and $GOPATH
 * Make sure to Add $GOROOT/bin and $GOPATH/bin to $PATH
* Install node.js
 * https://nodejs.org/download/
* Install go-bindata
 * run `go install github.com/jteeuwen/go-bindata`
* Install go-bindata-assetfs
 * run `go get github.com/elazarl/go-bindata-assetfs`
* In terminal enter `go get github.com/apokalyptik/fof/tools/team`
 * Might get an error, but ignore it for now
 * copy config.yaml.example to config.yaml. Add a slack API token (free teams at slack.com).
 * `cd $GOPATH/github.com/apokalyptik/fof/tools/team`
 * `make`
 * run `./team`

## Developing
Follow the above instructions. You will also need to install webpack binary (`npm install -g webpack`). Then, instead of `make` run `make dev`. This will recompile the react JS files as they are edited and saved. These changes can be viewed at http://localhost:8878
