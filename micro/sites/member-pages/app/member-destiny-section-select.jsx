var React = require('react')

module.exports = React.createClass({
	render: function() {
		return (
			<div>
				<div className="row">
					<div className="col-md-12">
						<h3>Select a section</h3>
					</div>
				</div>
				<div className="row">
					<div className="col-md-11 col-md-offset-1">
						<strong>PVP</strong> || PVE || Raid
						// TODO: Select
						// TODO: Part of route data...
					</div>
				</div>
			</div>
		);
	},
});

