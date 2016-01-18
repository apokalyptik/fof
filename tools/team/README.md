# Federation of Fathers Slack Integration Code


## Getting Started

**Install Go**
- https://golang.org/doc/install

**Install node.js**
- https://nodejs.org/download/
- If on OS X, Homebrew is recommended. If you don't have Homebrew:
  - `ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
  - Then run `brew install node`

**Setting up Go Environment**
- Run `mkdir -p $HOME/workspace/go`
- Set `$GOROOT` and `$GOPATH` in your `~/.bash_profile`
  - `export GOPATH=/Users/username/workspace/go` (change `username` to your username)
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

**Install go-uuid**
- run `go get github.com/pborman/uuid`

**Install team tool**
- run `go get github.com/apokalyptik/fof/tools/team` (*you might get an error, but ignore it for now*)

**Config and building**
- Copy `config.yaml.example` to `config.yaml`. 
- In your config, add a slack API token (*free teams at [slack.com](https://slack.com/)*) in `apiKey`.
- `cd $GOPATH/src/github.com/apokalyptik/fof/tools/team`
- run `make`
- run `./team`

**Logging In**
- In slack, type `/team` and copy the link from `fofbot`. Paste into your browser, and replace `team.fofgaming.com` with `http://localhost:8878` to login.


## Developing

Follow the above instructions. You will also need to install webpack binary (`npm install -g webpack`). Then, instead of `make` run `make && make dev`. This will recompile the react JS files as they are edited and saved. These changes can be viewed at [http://localhost:8878/dev](http://localhost:8878/dev). Once you run `make && make dev`, you will need to run `./team` in a new terminal window.

If you run into errors when running `make`, you may need to install [GNU command line tools](https://www.topbug.net/blog/2013/04/14/install-and-use-gnu-command-line-tools-in-mac-os-x/) and install [gnu-sed](http://stackoverflow.com/a/30005218).
