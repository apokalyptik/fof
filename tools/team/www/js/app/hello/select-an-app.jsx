React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');

module.exports = SelectAnApp = React.createClass({
	render: function() {
		var viewing = this.props.viewing;
		var lfgNow = "LFG Now";
		var lfgLater = "LFG Later";
		var lfgReport = "Report a Member";
		var buttonText = "";
		switch ( viewing ) {
			case "events":
				buttonText = lfgLater;
				break;
			case "lfg":
				buttonText = lfgNow;
				break;
			case "report":
				buttonText = lfgReport;
				break;
		}
		return (
			<div className="btn-group selectApp">
				<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
					<span className="glyphicon glyphicon-menu-hamburger" aria-hidden="true"></span><span className="hidden-xs">&nbsp;{buttonText}</span>
				</button>
				<ul className="dropdown-menu">
					<li><a href="#" data-value="events">{lfgLater}</a></li>
					<li><a href="#" data-value="lfg">{lfgNow}</a></li>
					<li><a href="#" data-value="report">{lfgReport}</a></li>
				</ul>
			</div>
		);
	},
	componentDidMount: function() {
		var viewing = this.props.viewing;
		$(".selectApp ul li a").bind("click tap",function(e){
			e.preventDefault();
			appName = $(e.target).data("value");
			Dispatcher.dispatch({ actionType: "set", key: "error", value: "" });
			Dispatcher.dispatch({ actionType: "set", key: "viewing", value: appName });

		});
		
	}
});

