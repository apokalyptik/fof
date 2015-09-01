var React = require('react');
var dispatcher = require('./dispatcher.js');

var DestinySectionSelect = React.createClass({
	change: function(e) {
		var routeData = this.props.route.data;
		routeData.section = e.target.value;
		dispatcher.dispatch({ type: "to", route: "member", params: routeData });
	},
	render: function() {
		return (
			<select value={this.props.route.data.section || ""} onChange={this.change}>
				<option value="">Select a section</option>
				<option value="pvp">PVP Stats</option>
			</select>
		);
	},
});

module.exports = React.createClass({
	change: function(e) {
		var routeData = this.props.route.data;
		routeData.game = e.target.value;
		dispatcher.dispatch({type: "to", route: "member", params: routeData})
	},
	render: function() {
		var subSelect = null;
		switch( this.props.route.data.game ) {
			case "destiny":
				subSelect = ( <DestinySectionSelect route={this.props.route}/> );
				break;
		}
		return (
			<div>
				<h3>
					<select onChange={this.change} value={this.props.route.data.game || ""}>
						<option value="">Select a Game...</option>
						<option value="destiny">Destiny</option>
					</select> {subSelect}
				</h3>
			</div>
		);
	},
});


