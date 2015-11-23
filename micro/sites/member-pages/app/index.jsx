// Detect and compile scss changes and include them (via JS)
// in the document head
require('./assets/style.scss');

// Add react as a window level object for the chrome react dev tools
// I am unsure if this is necessary. Just something I saw.
var React = require('react');
window.React = React;

// hashbang uri Routing...
//
// https://github.com/swipely/aviator
// defined route data will show up in this.state.routing
// for example the route:
//		"/:who { "/:what": {} }
// for the uri:
//		"#/foo/bar"
// will be available in <Root> with:
//		this.state.routing.who // == "foo"
//		this.state.routing.what // == "bar"
// 
// changes to the hashbang uri will automatically rerender.
//
// so long as you leave the target: {...} and "/*": "onChange",
// alone you do not *need* to specify a special handler for
// any route. for example:
//		"/some/:thing": {}
// will be handled properly because "/*" handles *all* changes.
//
// for webpack-dev-server this sees the <iframe> src="" parameter
// so for testing new routes where you don't have anything in the
// app setup yet you'll have to use the browsers developer tools
// to edit that directly
var routing = require('aviator');
var State = require('./state.js');
routing.setRoutes({
	target: {
		onChange: function( req ) {
			console.log( req.params );
			State.set( { routing: req.params } );
		},
	},
	"/*": "onChange",
});
routing.dispatch();

// Render the root element of our app. All else will flow from this.
var ReactDOM = require('react-dom');
var Root = require('./root.jsx');
ReactDOM.render( <Root/>, document.getElementById('root') );
