React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');

module.exports = React.createClass({
	click: function(e) {
		Dispatcher.dispatch({actionType: "set", key: "raid", value: this.props.data.uuid});
		e.preventDefault();
	},
	render: function() {
		className = "raid";
		if ( this.props.selected == this.props.data.uuid ) {
			className = className + " active";
		}
		var now = Date.now();
		var then = Date.parse(this.props.data.created_at);
		var ago = Math.round((now - then) / 8640000) / 10;
		if ( ago == 0 ) {
			ago = "0.0";
		}
		return (
			<div className={className}>
				<div className="row">
					<div className="col-md-9">
						<a onClick={this.click} href="#">{this.props.data.name}</a>
					</div>
					<div className="col-md-3">
						<em>{ago}</em> days ago | <em>{this.props.number}</em>
					</div>
				</div>
			</div>);
	}
});
