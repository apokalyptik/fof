React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');

module.exports = SelectAnApp = React.createClass({
	render: function() {
		var viewing = this.props.viewing;
		return (
			<select id="select-an-app" defaultValue={viewing} onChange={
				function(event) {
					Dispatcher.dispatch({
						actionType: "set",
						key: "viewing",
						value: event.target.value
					});
				}}>
				<option value="events">LFG Later</option>
				<option value="lfg">LFG Now</option>
			</select>
		);
	}
});

