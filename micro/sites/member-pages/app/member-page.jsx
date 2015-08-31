var React = require('react');
var MemberProfilePart = require('./member-profile.jsx');
var MemberGameSelect = require('./member-game-select.jsx');
var DestinyStats = require('./member-destiny-stats.jsx');

module.exports = React.createClass({
	render: function() {
		var wantMember = this.props.state.route.data.name.toLowerCase()
		var member = null;
		for ( var i=0; i<this.props.state.members.length; i++ ) {
			if ( this.props.state.members[i].username.toLowerCase() == wantMember ) {
				member = this.props.state.members[i];
				break;
			}
			if ( this.props.state.members[i].gamertag.toLowerCase() == wantMember ) {
				member = this.props.state.members[i];
				break;
			}
		}
		if ( member == null ) {
			window.setTimeout(function() {
				dispatcher.dispatch({type:"go", to: ""});
			}, 0);
			return((<div/>));
		}
	
		var gameSubDetails = null;
		switch ( this.props.state.route.data.section ) {
			case 'pvp':
				gameSubDetails = ( <DestinyStats state={this.props.state} member={member}/> );
				break;
		}

		return(
			<div className="container-fluid member">
				<div className="row">
					<div className="col-md-3">
						<div className="container-fluid profile-summary">
							<MemberProfilePart member={member} members={this.props.state.members} route={this.props.state.route}/>
						</div>
					</div>
					<div className="col-md-8">
						<div className="container-fluid">
							<MemberGameSelect route={this.props.state.route}/>
							{gameSubDetails}
						</div>
					</div>
				</div>
			</div>
		);
	}
});

