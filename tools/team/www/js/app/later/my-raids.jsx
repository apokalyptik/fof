React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Routing = require('aviator');

module.exports = React.createClass({
	isMember: function(raid) {
		for ( var i=0; i<raid.members.length; i++ ) {
			if ( raid.members[i] == this.props.state.username ) {
				return true;
			}
		}
		return false;
	},
	isAlt: function(raid) {
		if ( raid.alts == null ) {
			return false;
		}
		for ( var i=0; i<raid.alts.length; i++ ) {
			if ( raid.alts[i] == this.props.state.username ) {
				return true;
			}
		}
		return false;
	},
	raids: function() {
		var myRaids = [];
		for ( var c in this.props.state.raids ) {
			for ( var u in this.props.state.raids[c] ) {
				var r = this.props.state.raids[c][u];
				if ( this.isMember(r) || this.isAlt(r) ) {
					myRaids.push({raid: r, channel: c});
				}
			}
		}
		return myRaids;
	},
	select: function(e) {
		e.preventDefault();
		if ( e.target.value == "My Events" ) {
			Routing.navigate("/:section", { namedParams: { section: "later" } } );
		} else {
			var option = $(e.target).find('[value="'+e.target.value+'"]')[0];
			Routing.navigate(
				"/:section/:channel/:raid",
				{ namedParams:
					{ 
						section: "later",
						channel: option.dataset.channel,
						raid: option.dataset.uuid,
					}
				}
			);
		}
	},
	render: function() {
		var raids = this.raids();
		if ( raids.length < 1 ) {
			return(<span/>);
		}
		var raidlist = [<option key="none">My Events</option>];
		raids.forEach(function(entry) {
			raidlist.push(
				(<option
					value={entry.raid.uuid}
					key={entry.raid.uuid}
					data-channel={entry.channel}
					data-uuid={entry.raid.uuid}>
						{entry.raid.name}</option>)
			);
		})
		return (
			<div className="box col-xs-8 nopadding">
				<select className="form-control" value={this.props.routing.params.c} onChange={this.select}>
					{raidlist}
				</select>
			</div>
		);
	}
});
