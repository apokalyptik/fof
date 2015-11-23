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
		return (
			<div className="container-fluid">
				<div className="row">
					<div className="col-md-2 col-md-offset-10">TopNav</div>
				</div>
				<div className="row">
					<div className="col-md-2">
						<h1>logo</h1>

						<a>xbl link</a>
						<a>ping link</a>

						<h1>bio</h1>
						<ul>
							<li>city/state</li>
							<li>occupation</li>
							<li>gaming interest</li>
						</ul>

						<h1>stats</h1>
						<ul>
							<li><strong>Desting</strong></li>
							<li>time played</li>
							<li>pve kills</li>
							<li>raids completed</li>
							<li>pvp kills</li>
							<li>pvp k/d</li>
						</ul>
					</div>
					<div className="col-md-10">
						<h1>cheevos</h1>
						<h1>scheduled events</h1>
						<h1>channels</h1>
						<h1>xbox dvr</h1>
					</div>
				</div>
			</div>
		);
	}
});

