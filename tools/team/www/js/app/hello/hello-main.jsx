var React = require('react/addons');
var Dispatcher = require('../lib/dispatcher.jsx');
var Datastore = require('../lib/datastore.jsx');
var Routing = require('aviator');
var Config = require('../config.js');

var Hello = React.createClass({
	dispatch: function(event) {
		Routing.navigate("/:section", { namedParams: { section: event.target.value } });
		event.preventDefault();
	},
	render: function() {
		var sections = [];
		if ( Config.features.now ) {
			sections.push((<button value="now" onClick={this.dispatch}
				className="btn btn-block btn-default">Looking for Game Now</button>));
			sections.push(<span> or </span>);
		}
		sections.push((<button value="later" onClick={this.dispatch}
			className="btn btn-block btn-default">Looking for Game Later</button>));

		return(
			<div className="container fluid">
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>

				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						{sections}
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
				
				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<img style={{width:"100%"}} src="/logo.png"/>
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
				
				<div className="row"><div className="col-md-4 col-md-offset-4 center"><h6>
					As a member of the FoF Community and by using the team tools,
					you agree to follow the Federation of Fathers<br/> 
					<a target="_blank" href="http://fofgaming.com/fof/code-of-conduct/">Code of Conduct</a>
				</h6></div></div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
			</div>
		);
	}
});

module.exports = Hello;
