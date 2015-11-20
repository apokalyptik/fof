var webpack = require('webpack');
var path = require('path');
var HtmlwebpackPlugin = require('html-webpack-plugin');

var ROOT_PATH = path.resolve(__dirname);
var APP_PATH = path.resolve(ROOT_PATH, 'app');
var BUILD_PATH = path.resolve(ROOT_PATH, 'public_html');

var INDEX_FILE = path.resolve(BUILD_PATH, "index.html");
var INDEX_TMPL = path.resolve(APP_PATH, "assets/index-template.html");

module.exports = {
	entry: APP_PATH,
	output: {
		path: BUILD_PATH,
		filename: 'bundle.js'
	},
	devServer: {
		historyApiFallback: true,
		hot: true,
		inline: true,
		progress: true,
		// parse host and port from env so this is easy
		// to customize
		contentBase: BUILD_PATH,
		host: process.env.HOST,
		port: process.env.PORT || 8888
	},
	module: {
		loaders: [
			{ test: /\.scss$/, loaders: ["style", "css", "sass"], include: APP_PATH },
			{ test: /\.css$/, loaders: ['style', 'css'], include: APP_PATH },
			{ test: /\.jsx?$/,  exclude: /node_modules/, loader: "babel-loader" }
		]
	},
	resolve: {
		extensions: ['', '.js', '.jsx']
	},
	plugins: [
		new webpack.HotModuleReplacementPlugin(),
		new HtmlwebpackPlugin({
			title: 'FoF Member Pages',
			inject: true,
			template: "app/assets/index-template.html"
		})
	]
};
