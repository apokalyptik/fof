var Dispatcher = null;

function postRaid(what, data) {
	return jQuery.post("/rest/raid/" + what, data).fail(function( data ) {
		if ( data.status == 403 ) {
			location.reload( true );
		}
	});
}

var CountdownTimer = React.createClass({
	getInitialState: function() {
		return {
			secondsInitial: 0,
			secondsRemaining: 0,
			ticked: false,
		};
	},
	tick: function() {
		this.setState({secondsRemaining: this.state.secondsRemaining - 1, ticked: true});
		if (this.state.secondsRemaining <= 0) {
			clearInterval(this.interval);
		}
	},
	componentDidMount: function() {
		this.setState({
			secondsInitial: this.props.secondsRemaining,
			secondsRemaining: this.props.secondsRemaining
		});
		this.interval = setInterval(this.tick, 1000);
	},
	componentWillUnmount: function() {
		clearInterval(this.interval);
	},
	render: function() {
		var niceDisplay = "";
		var seconds = this.state.secondsRemaining;
		var minutes = Math.floor( seconds / 60 );
		seconds = seconds - ( 60 * minutes ); 
		if ( this.state.secondsRemaining >= 3600 ) {
			var hours = Math.floor(this.state.secondsRemaining/3600)
			minutes = minutes - ( hours * 60 );
			if ( minutes == 60 ) {
				minutes = 0;
				hours = hours + 1;
			}
			niceDisplay = hours + "h";
		}
		niceDisplay = niceDisplay + minutes + "m" + seconds + "s"
		var pct = 0;
		if ( this.state.ticked ) {
			if ( this.state.secondsRemaining == this.state.secondsInitial ) {
				pct = 100;
			} else {
				if ( this.state.secondsInitial != 0 ) {
					pct = 100 - Math.ceil(
						100 * ( this.state.secondsRemaining / this.state.secondsInitial )
					);
				}
			}
		}
		return (
			<div className="center">
				<div style={{textShadow: "0 0 1px #fff"}}>{niceDisplay}</div>
				<div className="progress-bar" style={{textAlign:"left", marginTop:"-1.7em"}}>
					<span className="center" style={{width: pct+"%"}}></span>
				</div>
			</div>
		);
	}
});

var Channel = React.createClass({
	select: function(e) {
		e.preventDefault();
		Dispatcher.dispatch({actionType: "set", key: "raid", value: ""});
		Dispatcher.dispatch({actionType: "set", key: "channel", value:this.props.name});
	},
	render: function() {
		var classes = [ "raidChannel" ];
		if ( this.props.name == this.props.selected ) {
			classes.push("active");
		}
		return(
			<div className={classes.join(" ")}>
				<a onClick={this.select} href="#">{this.props.name}</a> <span className="floatright">{this.props.number}</span>
			</div>
		);
	},
});

var ChannelList = React.createClass({
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

var Raid = React.createClass({
	click: function(e) {
		Dispatcher.dispatch({actionType: "set", key: "raid", value: this.props.data.uuid});
		e.preventDefault();
	},
	render: function() {
		className = "raid";
		if ( this.props.selected == this.props.data.uuid ) {
			className = className + " active";
		}
		var now = Date.now();
		var then = Date.parse(this.props.data.created_at);
		var ago = Math.round((now - then) / 8640000) / 10;
		if ( ago == 0 ) {
			ago = "0.0";
		}
		return (
			<div className={className}>
				<div className="row">
					<div className="col-md-9">
						<a onClick={this.click} href="#">{this.props.data.name}</a>
					</div>
					<div className="col-md-3">
						<em>{ago}</em> days ago | <em>{this.props.number}</em>
					</div>
				</div>
			</div>);
	}
});

var RaidList = React.createClass({
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
				if ( a.created_at < b.created_at ) {
					return -1;
				}
				if ( a.created_at > b.created_at ) {
					return 1;
				}
				return 0;
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

var AltMember = React.createClass({
	render: function() {
		if ( this.props.username != this.props.name ) {
			return (<div className="member alternate">@{this.props.name}</div>)
		}
		var leaveButton = (<span/>);
		if ( this.props.doLeaveButton ) {
			leaveButton = (
				<button 
					className="floatright btn btn-warning btn-xs" 
					onClick={this.props.leave} 
					href="#">leave</button>
			);
		}
		return (
			<div className="member alternate">
				<span className="me">@{this.props.name}</span>
				{leaveButton}
			</div>
		);
	}
});

var Member = React.createClass({
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

var MemberList = React.createClass({
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
		return(
			<div className="col-md-3">
				<h4>Members</h4>
				{myMemberList}
				{myAltList}
				<div style={{padding: "0.15em"}}>
					{joinBlock}
				</div>
			</div>
		);
	},
});

var HostForm = React.createClass({
	getInitialState: function() {
		return {
			error: "",
			channel: "",
			raid: "",
		}
	},
	submit: function(e) {
		if ( this.state.channel == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "please select a channel"});
			return;
		}
		if ( this.state.raid == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "please enter an event name"});
			return;
		}
		jQuery.post("/rest/raid/host", {channel: this.state.channel, raid: this.state.raid})
			.done(function(data) {
				Dispatcher.dispatch({actionType: "set", key: "error", value: ""});
				this.props.cancel(e)
			}.bind(this))
			.fail(function(data) {
				if ( data.status == 403 ) {
					location.reload(true);
				}
				Dispatcher.dispatch({actionType: "set", key: "error", value: responseText});
			}.bind(this));
	},
	handleRaid: function(event) { 
		this.setState({ "raid": event.target.value })
	},
	handleChannel: function(event) {
		this.setState({ "channel": event.target.value })
	},
	render: function() {
		var channels = [
			"",
		];
		for ( var i=0; i<this.props.channels.length; i++ ) {
			channels.push(this.props.channels[i]);
		}

		for ( var i=0; i<channels.length; i++ ) {
			if ( i == 0 ) {
				channels[i] = (<option key="" value="">-- Select a Channel --</option>);
			} else {
				channels[i] = (<option key={channels[i]} value={channels[i]}>{channels[i]}</option>);
			}
		}

		var errMsg = (<div/>);
		if ( this.state.error != "" ) {
			errMsg = ( <p className="bg-danger">{this.state.error}</p> );
		}
		return (
			<div className="col-md-6 col-md-offset-3">
			<h4>Host an Event</h4>
				<div className="form-group">
					<label htmlFor="channel">Channel to Host in</label>
					<select className="form-control" onChange={this.handleChannel}>
						{channels}
					</select>
				</div>
				<div className="form-group">
					<label htmlFor="name">Name of your Event</label>
					<input 
						onChange={this.handleRaid}
						type="text" className="form-control" id="name" placeholder="Event Name"/>
					<div className="row">
						<div className="col-md-8 col-md-offset-2">
							<em>Be sure to include the date, time, and time zone for your event</em>
						</div>
					</div>
				</div>
				{errMsg}
				<button
					onClick={this.submit}
					className="btn btn-primary">Host Event</button>
				<button
					onClick={this.props.cancel}
					className="btn btn-link">Never Mind...</button>
			</div>
		);
	}
});


var TeamApp = React.createClass({
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

var LFGSelectGame = React.createClass({
	render: function() {
		return(<div>Destiny</div>);
	}
});

var LFGAppLooking = React.createClass({
	cancel: function() {
		jQuery.post("/rest/lfg", { events: [], time: "0" })
			.done(function() {
				Dispatcher.dispatch({
					actionType: "lfg-looking",
					value: false
				});
			})
			.fail(function() {
				window.setTimeout(this.clear.bind(this), 500)
			});
	},
	ping: function(event) {
		event.preventDefault()
		var un = event.target.getAttribute('data-user');
		var ab = event.target.getAttribute('data-about');
		$.post('/rest/ping', { username: un, about: ab  })
			.done(function() {
				alert("Ping to " +un+ " successful")
			})
			.fail(function() {
				alert("Ping to " +un+ " failed" )
			})
	},
	renderSection: function(name) {
		var clearName = name.split(":").map(function(part) {
			return decodeURIComponent(part)
		}).join(" ");
		var peers = [];
		if ( typeof this.props.peers[name] != "undefined" ) {
			for ( var user in this.props.peers[name] ) {
				if ( user == this.props.username ) {
					continue;
				}
				var gt = this.props.peers[name][user].gamertag;
				var msg = "https://account.xbox.com/en-US/Messages?gamerTag=" + encodeURIComponent(gt)
				var pro = "https://account.xbox.com/en-us/profile?gamerTag=" + encodeURIComponent(gt)
				peers.push((
					<li key={name + "-" + user}>
						{gt}<br/>
						<a className="btn btn-default btn-xs" target="_blank" href={msg}>XBL Msg</a>&nbsp;
						<a className="btn btn-default btn-xs" target="_blank" href={pro}>XBL Profile</a>&nbsp;
						<a
							data-about={clearName}
							data-user={user}
							onClick={this.ping}
							className="btn btn-default btn-xs" 
							target="_blank" 
							href="#">Slack Ping</a>
					</li> ));
			}
		}
		return (
			<div key={"activity-"+name} className="col-md-3">
				<h5>{clearName}</h5>
				<ul className="lfg peers">{peers}</ul>
			</div>
		)
	},
	render: function() {
		var list = [];
		for ( var i in this.props.forWhat ) {
			if ( this.props.forWhat[i] == false ) {
				continue;
			}
			list.push(( this.renderSection(i) ) );
		}
	
		var gotRows = [];
		var wantRows = list.length / 4;
		for ( var i=0; i<wantRows; i++ ) {
			var thisRow = list.slice(i*4, (i+1)*4);
			if ( thisRow.length > 0 ) {
				gotRows.push(<div key={i} className="row">{thisRow}</div>);
			}
		}
		return (
			<div className="container-fluid">
				<div className="row">
					<div className="col-md-1">
						<button className="btn btn-default btn-block" onClick={this.cancel}>Reset</button>
						<br/>
						<CountdownTimer secondsRemaining={this.props.time * 60}/>
					</div>
					<div className="col-md-11">
						<div className="container-fluid">
								{gotRows}
						</div>
					</div>
				</div>
			</div>
		);
	}
});

var LFGApp = React.createClass({
	activities: [
		{
			name: "Prison of Elders",
			options: [
				{ name: "Level 28" },
				{ name: "Level 32" },
				{ name: "Level 34" },
				{ name: "Level 35" }
			]
		},
		{
			name: "Crota's End",
			options: [
				{ name: "Normal Mode" },
				{ name: "Hard Mode" }
			]
		},
		{
			name: "Vault of Glass",
			options: [
				{ name: "Normal Mode" },
				{ name: "Hard Mode" }
			]
		},
		{
			name: "Strikes",
			options: [
				{ name: "Weekly Nightfall" },
				{ name: "Weekly Heroic" },
				{ name: "Playlist"}
			],
		},
		{
			name: "Crucible",
			options: [
				{ name: "Trials of Osiris" },
				{ name: "Iron Banner" },
				{ name: "Daily PvP" },
				{ name: "Control" },
				{ name: "Skirmish" },
				{ name: "Clash" },
				{ name: "Rumble" },
				{ name: "Salvage" }
			],
		},
		{
			name: "Bounties",
			options: [
				{ name: "Queen" },
				{ name: "Eris" },
				{ name: "Daily Bounties" },
				{ name: "PvP Bounties" },
				{ name: "Exotic Bounty" }
			],
		},
		{
			name: "Story",
			options: [
			{ name: "Daily Mission" },
			{ name: "House of Wolves" },
			{ name: "The Dark Below" },
			{ name: "Earth" },
			{ name: "Moon" },
			{ name: "Venus" },
			{ name: "Mars" },
			],
		},
		{
			name: "Patrol",
			options: [
				{ name: "Public Events" },
				{ name: "Earth" },
				{ name: "Moon" },
				{ name: "Venus" },
				{ name: "Mars" }
			],
		},
		{
			name: "Farming",
			options: [
				{ name: "Glimmer" },
				{ name: "Patrol Missions" },
				{ name: "Ether Chest" }
			],
		},
	],
	isInMyOptions: function(option) {
		if ( typeof this.props.state.my[option] == "undefined" ) {
			return false;
		}
		if ( this.props.state.my[option] == true ) {
			return true;
		}
		return false;
	},
	isPersistedOption: function(option) {
		if ( typeof this.props.state.lfg[option] == "undefined" ) {
			return false;
		}
		if ( typeof this.props.state.lfg[option][this.props.state.username] != "undefined" ) {
			if ( typeof this.props.state.my[option] == "undefined" ) {
				Dispatcher.dispatch({
					actionType: "lfg",
					what: option,
					value: true
				});
				return true;
			}
		}
		return false;
	},
	isChecked: function(option) {
		if ( this.isInMyOptions(option) ) {
			return true;
		}
		if ( this.isPersistedOption(option) ) {
			return true;
		}
		return false;
	},
	check: function(event) {
		Dispatcher.dispatch({
			actionType: "lfg",
			what: event.target.value,
			value: event.target.checked
		});
	},
	time: function(event) {
		Dispatcher.dispatch({
			actionType: "lfg-time",
			value: event.target.value
		});
	},
	clear: function() {
		Dispatcher.dispatch({ actionType: "lfg-flush" });
	},
	submit: function() {
		var events = [];
		for ( var eventName in this.props.state.my ) {
			if ( this.props.state.my[eventName] ) {
				events.push( eventName );
			}
		}
		if ( events.length < 1 ) {
			retuen
		}
		jQuery.post("/rest/lfg", { events: events, time: this.props.state.time })
			.done(function() {
				Dispatcher.dispatch({
					actionType: "lfg-looking",
					value: true
				});
			})
			.fail(function() {
				window.setTimeout(this.submit.bind(this), 500)
			});
	},
	getLookers: function(name, dataset) {
		var lookers = ".";
		if ( typeof dataset[name] != "undefined" ) {
			lookers = 0;
			for ( var looker in dataset[name] ) {
				if ( looker == this.props.state.username ) {
					continue;
				}
				lookers = lookers + 1
			}
			if ( lookers == 0 ) {
				lookers = ".";
			}
		}
		return lookers;
	},
	render: function() {
		if ( this.props.state.looking == true ) {
			return (<LFGAppLooking
					activities={this.activities}
					prev={this.props.state.prevlfg}
					peers={this.props.state.lfg}
					username={this.props.state.username}
					time={this.props.state.time}
					forWhat={this.props.state.my} />);
		}
		var activities = this.activities;
		var activityBlock = [];
		var numChecked = 0;
		for ( var i in this.props.state.my ) {
			if ( this.props.state.my[i] ) {
				numChecked = numChecked + 1;
			}
		}
		for ( var i = 0; i<activities.length; i++ ) {
			var act = activities[i];
			var aname = encodeURIComponent(act.name);
			var opt = [];
			for ( var io = 0; io<act.options.length; io++ ) {
				var cname = aname + ":" + encodeURIComponent(act.options[io].name);
				var lookers = this.getLookers(cname, this.props.state.lfg);
				var oldlookers = this.getLookers(cname, this.props.state.prevlfg);
				var localClassName = "lfg count";
				var disabled = false;
				if ( !this.isChecked(cname) && numChecked >= 4 ) {
					if ( typeof this.props.state.my[cname] == "undefined" || !this.props.state.my[cname] ) {
						disabled = true;
					}
				}
				if ( lookers != oldlookers ) {
					localClassName = "lfg count updated";
				}
				opt.push( (<li key={"activity-"+i+"-"+io}><input
							value={cname}
							onChange={this.check}
							checked={this.isChecked(cname)}
							disabled={disabled}
							type="checkbox"/> {act.options[io].name}
							<span 
							className={localClassName}>{lookers}</span></li> ) );
			}
			activityBlock.push( (
				<div key={"activity-"+i+"-"+activityBlock.length} className="col-md-3">
					<ul className="lfgselect">
						<li><h4>{act.name}</h4><ul>{opt}</ul></li>
					</ul>
				</div>
			) );
		}
		var activityRows = [];
		var wantRows = activityBlock.length/4;
		for ( var i=0; i<wantRows; i++ ) {
			var thisRow = activityBlock.slice(i*4, (i+1)*4);
			if ( thisRow.length > 0 ) {
				activityRows.push(<div key={i} className="row">{thisRow}</div>);
			}
		}
		var actionWidgets = (
			<div>
				<br/>
				<select value={this.props.state.time} defaultValue="120" onChange={this.time}>
					<option value="30">30 minutes</option>
					<option value="60">1 hour</option>
					<option value="90">1 hour 30 minutes</option>
					<option value="120">2 hours</option>
				</select>
				<br/>
				<br/>
				<button onClick={this.submit}>Submit</button>
				&nbsp;
				<span className="greentext">or</span>
				&nbsp;
				<button onClick={this.clear}>Clear</button>
			</div>
		);
		return (
			<div className="container-fluid">
				<div className="row">
					<div className="col-md-2 center">
						<span className="greentext bold">
							Select up to 4 activities and set your play duration	
						</span>
						<br/>
						{actionWidgets}
					</div>
					<div className="col-md-10">
						<div className="container-fluid">
							{activityRows}
						</div>
					</div>
				</div>
				<div className="row">
					<div className="col-md-3 col-md-offset-5 center">
						<br/>
						{actionWidgets}
					</div>
				</div>
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
			</div>
		);
	}
});

var SelectAnApp = React.createClass({
	render: function() {
		var viewing = this.props.viewing;
		return (
			<select id="select-an-app" defaultValue={viewing} onChange={
				function(event) {
					Dispatcher.dispatch({
						actionType: "set",
						key: "viewing",
						value: event.target.value
					});
				}}>
				<option value="events">LFG Later</option>
				<option value="lfg">LFG Now</option>
			</select>
		);
	}
});

var App = React.createClass({
	getInitialState: function() {
		return Datastore.data
	},

	componentDidMount: function() {
		jQuery.getJSON("/rest/login/check")
			.done(function(data) {
				if ( typeof data.username == "string" && data.username != "" ) {
					Dispatcher.dispatch({actionType: "set", key: "cmd", value: data.cmd});
					Dispatcher.dispatch({actionType: "set", key: "checked", value: true});
					Dispatcher.dispatch({actionType: "set", key: "username", value: data.username});
					Dispatcher.dispatch({actionType: "username", value: data.username});
					Dispatcher.dispatch({actionType: "set", key: "authenticated", value: true});
				} else {
					Dispatcher.dispatch({actionType: "set", key: "cmd", value: data.cmd});
					Dispatcher.dispatch({actionType: "set", key: "checked", value: true});
				}
				Datastore.subscribe(this.acceptData)
				this.updateData();
			}.bind(this));
	},

	updateData: function() {
		if ( this.state.authenticated == false ) {
			window.setTimeout(this.updateData, 1000);
			return;
		}
		jQuery.getJSON("/rest/get?since="+this.state.updated_at)
			.done(function(data) {
				Dispatcher.dispatch({actionType: "serverStateUpdate", data: data});
			})
			.fail(function(data) {
				if ( data.status == 403 ) {
					location.reload(true);
				}
			})
			.always(function() {
				window.setTimeout(this.updateData, 250);
			}.bind(this))
	},

	acceptData: function(newData) {
		this.setState(newData);
	},

	render: function() {
		if ( this.state.authenticated == false ) {
			return(
				<div className="container-fluid">
					<div className="row">
						<div className="col-md-6 col-md-offset-3 center">
							<h2 className="dark">
							please use the slack command &ldquo;<strong>{this.state.command}</strong>&rdquo; to log in
							</h2>
						</div>
					</div>
				</div>
			);
		}

		if ( this.state.checked == false ) {
			return (<div/>);
		}

		if ( this.state.viewing == "hello" ) {
			return ( <Hello/> );
		}

		var crumbs = [
			( <li key="appselect" className="box">
				  <SelectAnApp key="selectanapp" viewing={this.state.viewing}/>
			  </li> )
		];

		var WorkSpace;
		switch ( this.state.viewing ) {
			case "events":
				WorkSpace = ( <TeamApp state={this.state}/> );
				break;
			case "lfg":
				WorkSpace = ( <LFGApp state={this.state.lfg}/> );
				crumbs.push( ( <li key="crumb-lfg" className="box"><LFGSelectGame/></li> ) );
				break;
		}

		return(
			<div>
				<div className="container-fluid nopadding">
					<div className="row nomargin">
						<div className="col-md-12 nomargin">
							<h2 className="nomargin">
								FoF @{this.state.username}
							</h2>
							<div id="crumb-bar">
								<ul className="breadcrumbs-lgr">
									{crumbs}
									<li className="rt"/>
								</ul>
							</div>
						</div>
					</div>
				</div>
				{WorkSpace}
			</div>
		);
	},
});

var Hello = React.createClass({
	dispatch: function(event) {
		Dispatcher.dispatch({
			actionType: "hello-choose",
			value: event.target.value}
		);
	},
	render: function() {
		return(
			<div className="container fluid">
				
				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<h2>Federation of Fathers</h2>
					</div>
				</div>

				<div className="row"><div className="col-md-1">&nbsp;</div></div>

				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<button value="lfg" onClick={this.dispatch}
							className="btn btn-block btn-default">Looking For Game Now</button>
						or
						<button value="events" onClick={this.dispatch}
							className="btn btn-block btn-default">Looking for Game Later</button>
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
				
				<div className="row">
					<div className="col-md-4 col-md-offset-4 center">
						<img style={{width:"100%"}} src="/logo.png"/>
					</div>
				</div>
				
				<div className="row"><div className="col-md-1">&nbsp;</div></div>
			</div>
		);
	}
});

Dispatcher = fluxify.dispatcher;

var LFGStore = {
	callbacks: [],
	data: {
		looking: false,
		my: {},
		lfg: {},
		prevlfg: {},
		time: "120",
		username: "",
	},
	since: "0",
	set: function(what, value) {
		this.data.my[what] = value;
		this.emitChange();
	},
	subscribe: function(callback) {
		this.callbacks.push(callback);
	},
	emitChange: function() {
		for( var i = 0; i < this.callbacks.length; i++ ) {
			this.callbacks[i]( this.data );
		}
	}
};

Dispatcher.register(function(payload) {
	switch( payload.actionType ) {
		case "username":
			LFGStore.data.username = payload.value;
			LFGStore.emitChange();
			break;
		case "lfg-flush":
			LFGStore.data = { my: {} };
			LFGStore.emitChange();
			break;
		case "lfg":
			LFGStore.set(payload.what, payload.value);
			break;
		case "lfg-time":
			LFGStore.data.time = payload.value;
			LFGStore.emitChange();
			break;
		case "lfg-looking":
			LFGStore.data.looking = payload.value;
			LFGStore.emitChange();
			break;
		case "lfg-lp-results":
			LFGStore.data.prevlfg = LFGStore.data.lfg;
			LFGStore.data.lfg = payload.value;
			LFGStore.emitChange();
			break;
		case "lfg-since":
			LFGStore.since = payload.value;
			break;
	}
});

function lfgLongPoll() {
	$.getJSON("/rest/lfg?since=" + LFGStore.since)
		.success(function(data) {
			Dispatcher.dispatch({actionType: "lfg-since", value: data.updated_at});
			Dispatcher.dispatch({actionType: "lfg-lp-results", value: data.lfg});
			window.setTimeout(lfgLongPoll, 250)
		})
		.fail(function() {
			window.setTimeout(lfgLongPoll, 2000)
		})
}
lfgLongPoll();

var Datastore = {
	callbacks: [],
	data: {
		raid: "", // Selected Raid UUID
		channel: "", // Selected Raid Channel
		authenticated: false,
		checkedUsername: false,
		username: "",
		checked: false,
		command: "",
		updated_at: "",
		hosting: false,
		channels: [],
		viewing: "hello",
		lfg: {}
	},
	setThing: function(thing, value) {
		this.data[thing] = value;
		this.emitChange();
	},
	subscribe: function(callback) {
		this.callbacks.push(callback);
	},
	emitChange: function() {
		for( var i = 0; i < this.callbacks.length; i++ ) {
			this.callbacks[i]( this.data );
		}
	}
}

Dispatcher.register(function(payload) {
	var doReRender = false;
	switch ( payload.actionType ) {
		case "hello-choose":
			Datastore.setThing( "viewing", payload.value );
			break;
		case "serverStateUpdate":
			for ( var i in payload.data ) {
				Datastore.data[i] = payload.data[i];
			}
			var channel = Datastore.data.channel;
			var raid = Datastore.data.raid;
			if ( channel != "" ) {
				if ( typeof payload.data.raids[channel] == "undefined" ) {
					Datastore.data.channel = "";
					Datastore.data.raid = "";
				} else {
					if ( raid != "" ) {
						if ( typeof Datastore.data.raids[channel][raid] == "undefined" ) {
							Datastore.data.raid = "";
						}
					}
				}
			}
			Datastore.emitChange();
		case "set":
			Datastore.setThing(payload.key, payload.value);
	}
});

var hash = {
	parts: {},
	parse: function() {
		var pieces = location.hash.substring(1).split("&");
		for ( var i=0; i<pieces.length; i++ ) {
			var part = pieces[i];
			var parts = part.split("=").map(function(bit) {
				return decodeURIComponent(bit)
			});
			if ( parts.length == 1 ) {
				parts.push("");
			}
			this.parts[parts[0]] = parts[1];
		}
	},
	get: function(bit) {
		if ( typeof this.parts[bit] == undefined ) {
			return null;
		}
		return this.parts[bit];
	},
	isset: function(bit) {
		if ( typeof this.parts[bit] == undefined ) {
			return false;
		}
		return true;
	}
}

LFGStore.subscribe(function(data) {
	data.username = Datastore.data.username;
	Datastore.setThing("lfg", data);
});

jQuery(document).ready(function() {
	hash.parse()
	switch ( hash.get("app") ) {
		case "lfgnow":
			Datastore.data.viewing = "lfg";
			break;
		case "lfglater":
			Datastore.data.viewing = "events";
			break;
	}
	React.render(<App />, document.getElementById('app'));
})
