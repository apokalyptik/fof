React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Datastore = require('../lib/datastore.jsx');
LFGAppLooking = require('./app-looking.jsx');

Dispatcher.register(function(payload) {
	switch ( payload.actionType ) {
		case "lfg-flush":
			for ( var i in Datastore.data.lfg.my ) {
				Datastore.data.lfg.my[i] = false;
			}
			for ( var e in Datastore.data.lfg.lfg ) {
				if ( typeof Datastore.data.lfg.lfg[e][Datastore.data.lfg.username] != "undefined" ) {
					Datastore.data.lfg.my[e] = false;
				}
			}
			Datastore.emitChange();
			break;
		case "lfg":
			Datastore.data.lfg.my[payload.what] = payload.value;
			Datastore.emitChange();
			break;
		case "lfg-time":
			Datastore.data.lfg.time = payload.value;
			Datastore.emitChange();
			break;
		case "lfg-looking":
			Datastore.data.lfg.looking = payload.value;
			Datastore.emitChange();
			break;
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
		{
			name: "Exploration",
			options: [
				{ name: "Gold Chest Hunting" },
				{ name: "Ghost Collecting" },
			]
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
	getMyEvents: function() {
		var username = this.props.state.username;
		var events = [];
		for ( var eventName in this.props.state.my ) {
			if ( this.props.state.my[eventName] ) {
				events.push( eventName );
			}
		}
		for ( var eventName in this.props.state.lfg ) {
			if ( typeof this.props.state.lfg[eventName][username] != "undefined" ) {
				if ( typeof this.props.state.my[eventName] == "undefined" ) {
					events.push( eventName );
				}
			}
		}
		return events;
	},
	submit: function() {
		var events = this.getMyEvents();
		if ( events.length < 1 ) {
			return
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
		var myEvents = this.getMyEvents();
		if ( this.props.state.looking == true ) {
			return (<LFGAppLooking
					activities={this.activities}
					prev={this.props.state.prevlfg}
					peers={this.props.state.lfg}
					username={this.props.state.username}
					time={this.props.state.time}
					forWhat={myEvents} />);
		}
		var activities = this.activities;
		var activityBlock = [];
		var numChecked = myEvents.length;
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
					disabled = true;
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
					<option value="1">1 minute</option>
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

module.exports = LFGApp;
