React = require('react/addons');

module.exports = React.createClass({
	render: function() {
		if ( this.props.username != this.props.name ) {
			return (<div className="member">@{this.props.name}</div>)
		}
		var leaveButton = (<span/>);
		if ( this.props.doLeaveButton ) {
			leaveButton = (
				<button className="floatright btn btn-warning btn-xs" onClick={this.props.leave} href="#">leave</button>
			);
		}
		return (
			<div className="member">
				<span className="me">@{this.props.name}</span>
				{leaveButton}
			</div>
		);
	}
});

