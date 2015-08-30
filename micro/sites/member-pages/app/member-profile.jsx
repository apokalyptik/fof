var React = require("react");

module.exports = React.createClass({
	render: function() {
		var gtPart = encodeURIComponent(this.props.member.gamertag);
		
		var messageXBL = "https://account.xbox.com/en-US/Messages?gamerTag=" + gtPart 
		var profileXBL = "https://account.xbox.com/en-us/profile?gamerTag=" + gtPart

		return (
			<div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					<h3>{this.props.member.gamertag}</h3>
				</div></div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					@{this.props.member.username}
				</div></div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					<br/>
					<img height="192" width="192" src={this.props.member.avatar}/>
					<br/>
				</div></div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					<a href={profileXBL} target="_blank">Xbox Profile</a>
				</div></div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					<a href={messageXBL} target="_blank">Xbox Message</a>
				</div></div>
			</div>
		);
	},
});
