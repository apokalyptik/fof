var React = require("react/addons");
var Route = require('route-parser');
var datastore = require('./datastore.js');
var dispatcher = require('./dispatcher.js');

var MemberProfilePart = require('./member-profile.jsx');
var MemberChooserPart = require('./member-chooser.jsx');
var MemberGameSelect = require('./member-game-select.jsx');
var DestinySelectPart = require('./member-destiny-section-select.jsx');
var DestinyStats = require('./member-destiny-stats.jsx');

function hashDidChange() {
	if ( window.location.hash.length < 2 ) {
		dispatcher.dispatch({type:"route", name: "", data: {}});
		return
	}
	var hash = window.location.hash.substring(1);
	for ( var i=0; i<routes.length; i++ ) {
		var m = routes[i].route.match(hash)
		if ( m ) {
			dispatcher.dispatch({type:"route", name: routes[i].name, data: m});
			return;
		}
	}
	window.location.replace(window.location.pathname + window.location.search + "#");
}

var MemberPage = React.createClass({
	render: function() {
		var wantMember = this.props.state.route.data.name.toLowerCase()
		var member = null;
		for ( var i=0; i<this.props.state.members.length; i++ ) {
			if ( this.props.state.members[i].username.toLowerCase() == wantMember ) {
				member = this.props.state.members[i];
				break;
			}
			if ( this.props.state.members[i].gamertag.toLowerCase() == wantMember ) {
				member = this.props.state.members[i];
				break;
			}
		}
		if ( member == null ) {
			window.setTimeout(function() {
				dispatcher.dispatch({type:"go", to: ""});
			}, 0);
			return((<div/>));
		}
		
		return(
			<div className="container-fluid member">
				<div className="row">
					<div className="col-md-3">
						<div className="container-fluid profile-summary">
							<MemberProfilePart member={member}/>
							<MemberChooserPart/>
						</div>
					</div>
					<div className="col-md-8">
						<div className="container-fluid">
			
							<MemberGameSelect/>
							<DestinySelectPart/>

							<div className="row">
								<div className="col-md-12">
									<h3>Details Below...</h3>
								</div>
							</div>

							<DestinyStats state={this.props.state} member={member}/>

						</div>
					</div>
				</div>
			</div>
		);
	}
});

var App = React.createClass({
	render: function() {
		if ( this.state.loaded != "v1" ) {
			return (<div>Loading Data...</div>);
		}
		if ( this.state.route.name == "member" ) {
			return ( <MemberPage state={this.state}/> );
		}
		return (<h1>{JSON.stringify(this.state.route)}</h1>);
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
				dispatcher.dispatch({type: "set", key: "lb.pvp", val: data});
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
	{ name: "lboard",	route: new Route('s/:stat(/:statRatioBy(/:statRatioOver))') },
];

$(document).ready(function() {
	$(window).on('hashchange', hashDidChange)
	hashDidChange()
	React.render( <App/>, document.getElementById('app') );
});
