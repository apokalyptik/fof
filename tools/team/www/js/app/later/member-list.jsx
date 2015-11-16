React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');

Member = require('./member.jsx');
AltMember = require('./alt-member.jsx');

function postRaid(what, data) {
	return jQuery.post("/rest/raid/" + what, data).fail(function( data ) {
		if ( data.status == 403 ) {
			location.reload( true );
		}
	});
}

module.exports = React.createClass({
	raidName: function() {
		return this.props.data[this.props.channel][this.props.raid].name;
	},
	raidPostData: function() {
			return { channel: this.props.channel, raid: this.raidName() };
	},
	join:     function() { postRaid( "join", this.raidPostData() )      },
	joinAlt:  function() { postRaid( "join-alt", this.raidPostData() )  },
	leave:    function() { postRaid( "leave", this.raidPostData() )     },
	leaveAlt: function() { postRaid( "leave-alt", this.raidPostData() ) },
	ping:     function() { postRaid( "ping", this.raidPostData() )      },
	finish:   function() { postRaid( "finish",this.raidPostData() )     },
	render: function() {
		var myMemberList = (
			<strong>
				Please select a raid to see the member list and be able to join or part
			</strong>
		);
		var myAltList = (<div/>);
		var isMember = false;
		if ( this.props.channel != "" && typeof this.props.data[this.props.channel] != "undefined" ) {
			if ( this.props.raid != "" && typeof this.props.data[this.props.channel][this.props.raid] != "undefined" ) {
				memberList = this.props.data[this.props.channel][this.props.raid].members;
				if ( memberList.length < 1 ) {
					myMemberList = (<span>This raid has no members</span>);
				} else {
					myMemberList = []
					var lastSelf = -1;
					for ( var i = 0; i<memberList.length; i++ ) {
						if ( memberList[i] == this.props.username ) {
							lastSelf = i;
						}
					}
					for ( var i = 0; i<memberList.length; i++ ) {
						if ( memberList[i] == this.props.username ) {
							isMember = true;
						}
						var doLeaveButton = false;
						if ( i == lastSelf ) {
							doLeaveButton = true;
						}
						myMemberList[i] = (
							<Member
								channel={this.props.channel}
								raid={this.props.data[this.props.channel][this.props.raid].name}
								key={this.props.raid.uuid + "-" + memberList[i] + "-" + i}
								name={memberList[i]}
								username={this.props.username}
								leader={this.props.data[this.props.channel][this.props.raid].members[0]}
								doLeaveButton={doLeaveButton}
								leave={this.leave}
								finish={this.props.finish}/>
						);
					}
				}
				var altList = this.props.data[this.props.channel][this.props.raid].alts;
				myAltList = [(
					<h4 key="alt" className="alternate">Alternates</h4>)];
				if ( typeof altList == "object" && altList != null && altList.length > 0 ) {
					var lastSelf = -1;
					for ( var i = 0; i<altList.length; i++ ) {
						if ( altList[i] == this.props.username ) {
							lastSelf = i;
						}
					}
					for ( i=0; i<altList.length; i++ ) {
						var doLeaveButton = false;
						if ( i == lastSelf ) {
							doLeaveButton = true;
						}
						myAltList.push(
							<AltMember
								channel={this.props.channel}
								raid={this.props.data[this.props.channel][this.props.raid].name}
								key={this.props.raid.uuid + "-alt-" + altList[i] + "-" + i}
								name={altList[i]}
								username={this.props.username}
								leader={this.props.data[this.props.channel][this.props.raid].members[0]}
								doLeaveButton={doLeaveButton}
								leave={this.leaveAlt}
								finish={this.props.finish}/>
						);
					}
				}

				var btnJ  = ( <button className="btn btn-success" onClick={this.join}>join</button> );
				var btnJA = ( <button className="btn btn-success" onClick={this.joinAlt}>join-alt</button> );
				var btnP  = ( <button className="btn btn-warning" onClick={this.ping} href="#">ping</button> );
				var btnF  = ( <button className="btn btn-danger" onClick={this.finish} href="#">finish</button> );

				var joinBlock = ( <div>{btnJ}&nbsp;{btnJA}</div> );

				isAdmin = false;
				for ( var i=0; i<this.props.admins.length; i++ ) {
					if ( this.props.admins[i] == this.props.username ) {
						isAdmin = true;
						break;
					}
				}

				if ( isMember || isAdmin ) {
					var leader = this.props.data[this.props.channel][this.props.raid].members[0]
					if ( leader == this.props.username || isAdmin ) {
						joinBlock = ( <div> {btnJ}&nbsp;{btnJA}&nbsp;{btnP}&nbsp;{btnF}</div> );
					}
				}
			}
		}
		
		var ics = null;
		if ( this.props.raid !== "" ) {
			ics = (<span style={{"float":"right"}}><a href="#" onClick={function(e) {
				var start = new Date(this.raid_time);
				var stop = new Date(new Date(this.raid_time).getTime() + 1800000);
				var cal = window.ics();
				cal.addEvent(this.raid_title, this.raid_title, "Federation of Fathers", start, stop);
				console.log(start);
				console.log(stop);
				console.log(cal);
				cal.download(this.raid_title)
				e.stopPropagation()
				e.preventDefault()
			}.bind(this.props.data[this.props.channel][this.props.raid])}>📅</a></span>)
		}

		return(
			<div className="col-md-3">
			<h4>Members{ics}</h4>
				{myMemberList}
				{myAltList}
				<div style={{padding: "0.15em"}}>
					{joinBlock}
				</div>
			</div>
		);
	},
});
