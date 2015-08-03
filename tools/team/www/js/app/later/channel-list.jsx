React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Datastore = require('../lib/datastore.jsx');
Channel = require('./channel.jsx');

module.exports = React.createClass({
	render: function() {
		var channelList = (<div/>)
		if ( typeof this.props.data != "undefined" ) {
			channelList = (<strong>There are currently no raids being hosted</strong>);
		}
		var channels = [];
		var channelRaids = [];
		for ( var c in this.props.data ) {
			var raids = 0;
			for ( var r in this.props.data[c] ) {
				raids = raids + 1;
			}
			if ( raids > 0 ) {
				channels.push(c);
				channelRaids.push(raids);
			}
		}
		if ( channels.length > 0 ) {
			for ( var i=0; i<channels.length; i++ ) {
				channels[i] = ( <Channel
					key={channels[i]}
					number={channelRaids[i]}
					name={channels[i]}
					selected={this.props.selected}/> );
			}
			channelList = channels
		}
		return(
			<div className="col-md-3">
				<h4>Channels</h4>
				{channelList}
				{this.props.host}
			</div>
		);
	},
});

