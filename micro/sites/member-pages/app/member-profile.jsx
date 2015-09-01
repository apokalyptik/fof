var React = require("react");
var dispatcher = require('./dispatcher.js');

var Chooser = React.createClass({
	change: function(e) {
		var routeData = this.props.route.data;
		routeData.name = e.target.value;
		dispatcher.dispatch({ type: "to", route: "member", params: routeData });
	},
	render: function() {
		var options = [];
		for( var i=0; i<this.props.members.length; i++ ) {
			var member = this.props.members[i];
			options.push((
				<option key={member.username} value={member.username}>{member.gamertag}</option>
			));
		}
		return(
			<select value={this.props.member.username} onChange={this.change}>
				{options}
			</select>
		);
	}
});

module.exports = React.createClass({
	render: function() {
		var gtPart = encodeURIComponent(this.props.member.gamertag);
		
		var messageXBL = "https://account.xbox.com/en-US/Messages?gamerTag=" + gtPart 
		var profileXBL = "https://account.xbox.com/en-us/profile?gamerTag=" + gtPart
		
		return (
			<div>
				<div className="row"><div className="col-md-10 col-md-offset-1">
					<h3>
						<Chooser member={this.props.member} members={this.props.members} route={this.props.route}/>
					</h3>
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
