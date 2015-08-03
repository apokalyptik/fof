var React = require('react/addons');
var Dispatcher = require('../lib/dispatcher.jsx');
var Datastore = require('../lib/datastore.jsx');

Dispatcher.register(function(payload) {
	var doReRender = false;
	switch ( payload.actionType ) {
		case "hello-choose":
			Datastore.setThing( "viewing", payload.value );
			break;
	}
});

var Hello = React.createClass({
	dispatch: function(event) {
		Dispatcher.dispatch({
			actionType: "hello-choose",
			value: event.target.value}
		);
	},
	render: function() {
		return(
			<div className="container fluid">
				
				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<h2>Federation of Fathers</h2>
					</div>
				</div>

				<div className="row"><div className="col-md-1">&nbsp;</div></div>

				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<button value="lfg" onClick={this.dispatch}
							className="btn btn-block btn-default">Looking for Game Now</button>
						or
						<button value="events" onClick={this.dispatch}
							className="btn btn-block btn-default">Looking for Game Later</button>
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
				
				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<img style={{width:"100%"}} src="/logo.png"/>
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
			</div>
		);
	}
});

module.exports = Hello;
