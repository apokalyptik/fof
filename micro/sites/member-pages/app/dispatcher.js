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
		case "to":
			for ( i=0; i<routes.length; i++ ) {
				if ( routes[i].name == d.route ) {
					window.location.href = window.location.pathname + window.location.search + '#' + routes[i].route.reverse(d.params);
					return
				}
			}
			break;
		case "go":
			window.location.href = window.location.pathname + window.location.search + '#' + d.to;
			break;
	}
});

module.exports = dispatcher
