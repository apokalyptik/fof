React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Raid = require('./raid.jsx');

module.exports = React.createClass({
	render: function() {
		var raidList = (
			<strong>
				Please select a channel to see the raid list
			</strong>
		);
		var channel = this.props.channel;
		if ( channel != "" ) {
			var raids = [];
			for ( var uuid in this.props.data[channel] ) {
				raids.push( this.props.data[channel][uuid] );
			}
			raids.sort(function(a, b) {
				var aDate = new Date(a["raid_time"]);
				var bDate = new Date(b["raid_time"]);

				return aDate.getTime() - bDate.getTime();
			});
			raidList = [];
			for ( var i=0; i<raids.length; i++ ) {
				var raid = raids[i];
				raidList.push( (<Raid
					key={raid.uuid}
					leader={raids[i].members[0]}
					select={this.props.select}
					selected={this.props.selected}
					number={raids[i].members.length}
					data={raid}/>) );
			}
			if ( raidList.length < 1 ) {
				raidList = ( <span>This channel has no raids</span> );
			}
		}
		return(
			<div className="col-md-6">
				<h4>Events</h4>
				{raidList}
			</div>
		);
	},
});

