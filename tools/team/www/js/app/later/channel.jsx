React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Datastore = require('../lib/datastore.jsx');

module.exports = React.createClass({
	select: function(e) {
		e.preventDefault();
		Dispatcher.dispatch({actionType: "set", key: "raid", value: ""});
		Dispatcher.dispatch({actionType: "set", key: "channel", value:this.props.name});
	},
	render: function() {
		var classes = [ "raidChannel" ];
		if ( this.props.name == this.props.selected ) {
			classes.push("active");
		}
		return(
			<div className={classes.join(" ")}>
				<a onClick={this.select} href="#">{this.props.name}</a> <span className="floatright">{this.props.number}</span>
			</div>
		);
	},
});

