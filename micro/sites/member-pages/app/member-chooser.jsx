var React = require('react')

module.exports = React.createClass({
	render: function() {
		return(
			<div className="row">
				<div className="col-md-10 col-md-offset-1">
					<h3>Member Chooser</h3>
					<input type="text" className="form-control" placeholder="find another user"/>
					<input type="submit" className="form-control" value="View This Member"/>
					// TODO: Select
				</div>
			</div>
		);
	},
});
