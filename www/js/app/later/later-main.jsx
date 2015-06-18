React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Datastore = require('../lib/datastore.jsx');

ChannelList = require('./channel-list.jsx');
RaidList = require('./raid-list.jsx');
MemberList = require('./member-list.jsx');
HostForm = require('./host-form.jsx');

module.exports = React.createClass({
	render: function() {

		if ( typeof this.props.state.raids == "undefined" ) {
			return (<div/>);
		}

		if ( this.props.state.hosting ) {
			return (
				<div className="container-fluid">
					<div className="row">
						<HostForm channels={this.props.state.channels} cancel={function() {
							Dispatcher.dispatch({actionType: "set", key: "hosting", value: false});
						}.bind(this)}/>
					</div>
				</div>
			);
		}
		
		var hostButton = (
			<button
				onClick={function(e) {
					e.preventDefault();
					Dispatcher.dispatch({actionType: "set", key: "hosting", value: true});	
				}.bind(this)}
				className="btn btn-default btn-block btn-success">Host an Event</button>
		);

		return (
			<div className="container-fluid">
				<div className="row">
					<ChannelList
						data={this.props.state.raids}
						selected={this.props.state.channel}
						host={hostButton}/>
					<RaidList data={this.props.state.raids}
						channel={this.props.state.channel}
						selected={this.props.state.raid}/>
					<MemberList
						username={this.props.state.username}
						channel={this.props.state.channel}
						raid={this.props.state.raid}
						data={this.props.state.raids}
						admins={this.props.state.admins}/>
				</div>
			</div>
		);
	}
});

