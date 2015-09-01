var Flux = require('flux');
var datastore = require('./datastore.js');

var dispatcher = new Flux.Dispatcher();

dispatcher.register(function(d) {
	switch( d.type ) {
		case "set":
			datastore.set(d.key, d.val);
			break;
		case "route":
			datastore.set("route", { name: d.name, data: d.data });
			break;
		case "go":
			window.location.href = window.location.pathname + window.location.search + '#' + d.to;
			break;
	}
});

module.exports = dispatcher
