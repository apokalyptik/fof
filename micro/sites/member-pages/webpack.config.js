var webpack = require("webpack");

var production = JSON.parse(process.env.PROD_DEV || "0");

module.exports = {
	entry: "./app/main.jsx",
	output: {
		filename: "bundle.js"
	},
	module: {
		loaders: [
			{ test: /\.jsx$/, loader: "jsx-loader?insertPragma=React.DOM&harmony" },
			{ test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"},
			{ test: /\.scss$/, loader: 'style!css!sass' },
			// { test: /\.css$/, loader: "style-loader!css-loader" },
			// { test: /\.png$/, loader: "url-loader?limit=100000" },
			// { test: /\.jpg$/, loader: "file-loader" }
		]
	},
	plugins: []
}

if ( production ) {
	module.exports.plugins.push(new webpack.optimize.UglifyJsPlugin({minimize: true}));
}

if ( !production ) {
	module.exports.devtool = "source-map";
}

