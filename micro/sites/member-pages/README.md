initializing (basically pull in a bunch of libs and setup some tooling...): `npm install`

dev server: `npm run dev` (or `make dev`) and visit http://localhost:8888/webpack-dev-server/bundle

build for production (builds into public_html): `npm run build` (or `make`)

build with dev stuff (source maps, etc): `npm run dev-build` (or `make dev-build`)

Most of the setup is done in `app/index.jsx`

The entry point into the whole app is `app/root.jsx`
