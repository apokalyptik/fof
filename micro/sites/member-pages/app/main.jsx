var React = require("react/addons");
var Route = require('route-parser');
var datastore = require('./datastore.js');
var dispatcher = require('./dispatcher.js');

var MemberPage = require('./member-page.jsx');
var LeaderBoard = require('./leaderboard.jsx');

var ComputePVP = require('./computed-pvp-stats.js');

function hashDidChange() {
	if ( window.location.hash.length < 2 ) {
		dispatcher.dispatch({type:"route", name: "", data: {}});
		return
	}
	var hash = window.location.hash.substring(1);
	if ( hash.substring(hash.length - 1) == "/" ) {
		dispatcher.dispatch({type:"go", to: hash.substring(0, hash.length - 1)});
		return;
	}
	for ( var i=0; i<routes.length; i++ ) {
		var m = routes[i].route.match(hash)
		if ( m ) {
			dispatcher.dispatch({type:"route", name: routes[i].name, data: m});
			return;
		}
	}
	window.location.replace(window.location.pathname + window.location.search + "#");
}

var App = React.createClass({
	render: function() {
		if ( this.state.loaded != "v1" ) {
			return (<div>Loading Data...</div>);
		}
		if ( this.state.route.name == "member" ) {
			return ( <MemberPage state={this.state}/> );
		}
		if ( this.state.route.name == "lboard" ) {
			return ( <LeaderBoard state={this.state}/> );
		}
		return (<h1>404 Not Found</h1>);
	},
	getInitialState: function() {
		return datastore.data;
	},
	componentDidMount: function() {
		datastore.listen(this.setState.bind(this));
		window.setTimeout(this.getUserJSON, 1000);
	},
	getUserJSON: function() {
		$.getJSON( "http://fofgaming.com:8880/fof/members.json" )
			.done(function(data) {
				dispatcher.dispatch({type: "set", key: "members", val: data});
				window.setTimeout(this.getLeaderboardPVP, 1000);
			}.bind(this))
			.fail(function(data) {
				window.setTimeout(this.getUserJSON, 1000);
			}.bind(this));
	},
	getLeaderboardPVP: function() {
		$.getJSON( "http://fofgaming.com:8880/destiny/stats/leaderboard/pvp.json" )
			.done(function(data) {
				dispatcher.dispatch({type: "set", key: "lb.pvp", val: ComputePVP(data)});
				dispatcher.dispatch({type: "set", key: "loaded", val: "v1"});
			}.bind(this))
			.fail(function(data) {
				window.setTimeout(this.getLeaderboardPVP, 1000);
				this.getLeaderboardPVP();
			}.bind(this));
	},
});

var routes = [
	{ name: "member",	route: new Route('m/:name(/:game(/:section))') },
	{ name: "lboard",	route: new Route('s/:type/:stat') },
];

dispatcher.register(function(d) {
	switch( d.type ) {
		case "to":
			for ( var i=0; i<routes.length; i++ ) {
				if ( routes[i].name == d.route ) {
					window.location.href = window.location.pathname + window.location.search + '#' + routes[i].route.reverse(d.params);
					return
				}
			}
			break;
	}
});

$(document).ready(function() {
	$(window).on('hashchange', hashDidChange)
	hashDidChange()
	React.render( <App/>, document.getElementById('app') );
});
