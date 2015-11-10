# Federation of Fathers Slack Raid Integration Code


## Getting Team Tool Running

**Install Go**
- https://golang.org/doc/install

**Install node.js**
- https://nodejs.org/download/

**Setting up Go Environment**
- In terminal, run `mkdir $HOME/go`
- Set $GOROOT and $GOPATH in your `~/.bash_profile`
  - `export GOPATH=/Users/username/go` (change `username` to your username)
- Add $GOROOT/bin and $GOPATH/bin to $PATH
  - `export PATH=$PATH:$GOROOT/bin:$GOPATH/bin`
 
**Install go-bindata**
- run `go get github.com/jteeuwen/go-bindata`
- `cd $GOPATH/src/github.com/jteeuwen/go-bindata`
- run `make`

**Install go-bindata-assetfs**
- run `go get github.com/elazarl/go-bindata-assetfs`
- `cd $GOPATH/src/github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs`
- `go build`
- `go install`

**Config and building**
- In terminal enter `go get github.com/apokalyptik/fof/tools/team` (*you might get an error, but ignore it for now*).
- Copy `config.yaml.example` to `config.yaml`. Add a slack API token (*free teams at [slack.com](https://slack.com/)*).
- `cd $GOPATH/src/github.com/apokalyptik/fof/tools/team`
- run `make`
- run `./team`

**Logging In**
- In slack, type `/team` and copy the link from `fofbot`. Paste into your browser, and replace `team.fofgaming.com` with `http://localhost:8878` to login.


## Developing

Follow the above instructions. You will also need to install webpack binary (`npm install -g webpack`). Then, instead of `make` run `make && make dev`. This will recompile the react JS files as they are edited and saved. These changes can be viewed at [http://localhost:8878/dev](http://localhost:8878/dev). 

If you run into errors when running `make` or `make && make dev`, you may also need to install browserify (`npm install -g browserify`). The first time you run `make` you may need to use `sudo` (for Mac). After the initial build, you should be able to run `make` without `sudo`.
