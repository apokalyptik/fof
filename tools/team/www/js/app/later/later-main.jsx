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
						routing={this.props.routing}
						data={this.props.state.raids}
						selected={this.props.routing.params.b}
						host={hostButton}/>
					<RaidList data={this.props.state.raids}
						routing={this.props.routing}
						channel={this.props.routing.params.b}
						selected={this.props.routing.params.c}/>
					<MemberList
						routing={this.props.routing}
						username={this.props.state.username}
						channel={this.props.routing.params.b}
						raid={this.props.routing.params.c}
						data={this.props.state.raids}
						admins={this.props.state.admins}/>
				</div>
			</div>
		);
	}
});

