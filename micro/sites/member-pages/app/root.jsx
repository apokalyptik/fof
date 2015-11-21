var React = require('react');
var State = require('./state.js');

module.exports =  React.createClass({
	/**
	 * Executes exactly once on the client. Unless the object is destroyed and a new one created...
	 */
	componentWillMount: function() {
		// Set our state to the global state
		this.setState( State.get().toObject() );
		// Subscribe to any and all global state changes for re-rendering purposes
		State.register(function(data) { this.setState( data.toObject() ) }.bind(this))
	},
	render: function() {
		return (<h1>Nothing yet...</h1>);
	}
});

