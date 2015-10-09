var webpack = require("webpack");
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var production = JSON.parse(process.env.PRODUCTION || "0");

module.exports = {
	entry: "./www/js/app/main.jsx",
	output: {
		filename: ""
	},
	module: {
		loaders: [
			{ test: /\.jsx$/, loader: "jsx-loader?insertPragma=React.DOM&harmony" },
			{ test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"},
			{ test: /\.scss$/, loader: ExtractTextPlugin.extract('style-loader', 'css-loader!sass-loader'), },
		]
	},
	plugins: [
		new ExtractTextPlugin("./www/css/style.css"),
	]
}

if ( production ) {
	module.exports.plugins.push(new webpack.optimize.UglifyJsPlugin({minimize: true, compress: { warnings: false } }));
	module.exports.output.filename = "./www/js/production.js";
}

if ( !production ) {
	module.exports.devtool = "source-map";
	module.exports.output.filename = "./www/js/development.js";
}
