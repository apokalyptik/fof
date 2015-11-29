var React = require('react/addons');
var Dispatcher = require('../lib/dispatcher.jsx');
var Datastore = require('../lib/datastore.jsx');
var Routing = require('aviator');

var Hello = React.createClass({
	dispatch: function(event) {
		Routing.navigate("/:section", { namedParams: { section: event.target.value } });
		event.preventDefault();
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
						<button value="now" onClick={this.dispatch}
							className="btn btn-block btn-default">Looking for Game Now</button>
						or
						<button value="later" onClick={this.dispatch}
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
