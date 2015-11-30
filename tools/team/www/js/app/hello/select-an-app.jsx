React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Config = require('../config.js');
Routing = require('aviator');

module.exports = SelectAnApp = React.createClass({
	render: function() {
		var lfgNow = "LFG Now";
		var lfgLater = "LFG Later";
		var lfgReport = "Submit a CoC Claim";
		var buttonText = "";
		switch ( this.props.routing.params.a ) {
			case "later":
				buttonText = lfgLater;
				break;
			case "now":
				buttonText = lfgNow;
				break;
			case "report":
				buttonText = lfgReport;
				break;
		}
		var items = [
				( <li><a href="#" data-value="later">{lfgLater}</a></li> ),
		];
		if ( Config.features.now ) {
			items.push(( <li><a href="#" data-value="now">{lfgNow}</a></li> ));
		}
		if ( Config.features.report ) {
			items.push( ( <li><a href="#" data-value="report">{lfgReport}</a></li> ) );
		}
		return (
			<div className="btn-group selectApp">
				<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
					<span className="glyphicon glyphicon-menu-hamburger" aria-hidden="true"></span><span className="hidden-xs">&nbsp;{buttonText}</span>
				</button>
				<ul className="dropdown-menu">{items}</ul>
			</div>
		);
	},
	componentDidMount: function() {
		$(".selectApp ul li a").bind("click tap",function(e){
			e.preventDefault();
			appName = $(e.target).data("value");
			Dispatcher.dispatch({ actionType: "set", key: "error", value: "" });
			Routing.navigate('/:section', { namedParams: { section: appName } } );
		});
		
	}
});

