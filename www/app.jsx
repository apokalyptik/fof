var active = { authenticated: false, channel: null, raid: null };

var returnFalse = function( e ) { e.stopPropagation(); return false; }

var Channel = React.createClass({
	render: function() {
		var classes = [ "raidChannel" ];
		if ( this.props.name == active.channel ) {
			classes.push("active");
		}
		return(
			<div className={classes.join(" ")}>
				<a href="#" onClick={this.props.select(this.props.name)}>{this.props.name}</a>
			</div>
		);
	},
});

var ChannelList = React.createClass({
	render: function() {
		var channelList = (
			<strong>There are currently no raids being hosted</strong>
		);
		var channels = [];
		for ( var c in this.props.data ) {
			if ( c[0] != "@" ) {
				channels.push(c)
			}
		}
		if ( channels.length > 0 ) {
			var channelList = channels.map(function (channelName) {
				return ( <Channel name={channelName} select={this.props.select}/> );
			});
		}
		return(
			<div className="col-md-3">
				<h2>Channels</h2>
				{channelList}
			</div>
		);
	},
});

var RaidList = React.createClass({
	render: function() {
		var raidList = (
			<strong>
				Please select a channel to see the raid list
			</strong>
		);
		return(
			<div className="col-md-5">
				<h2>Raids</h2>
				{raidList}
			</div>
		);
	},
});

var MemberList = React.createClass({
	render: function() {
		var memberList = (
			<strong>
				Please select a raid to see the member list and be able to join or part
			</strong>
		);
		return(
			<div className="col-md-4">
				<h2>Members</h2>
				{memberList}
			</div>
		);
	},
});

var LoginInit = React.createClass({
	render: function() {
		return(
			<div className="col-md-12">
				<form>
					<strong>Please enter your slack username or email address</strong>
					<input type="text" value="">
				</form>
			</div>
		);
	}
});

var Login = React.createClass({
	getIntialState: function() {
		return {
			username: "",
			authenticated: false,
			step: 0,
		}
	},
	render: function() {
		if ( self.state.step == 0 ) {
			return (<div><LoginInit/></div>);
		}
	},
});

var App = React.createClass({
	getIntialState: function() {
		return {};
	},
	selectChannel: function(name) {
		active.channel = name;
		this.setState(this.state)
	},
	render: function() {
		if ( active.authenticated == false ) {
			return( <div><Login/></div> );
		}
		var data = {
			"@active": active,
			"destiny-raid-crota": {
				"aaa": {
					"name": "cheese-tastic!",
					"members": ["demitriousk"],
				},
				"bbb": {
					name:"foo",
					members:[],
				},
			},
			"destiny-raid-vog": {
				"ccc": {},
			},
		}
		return(
			<div className="container-fluid">
				<div className="row">
					<ChannelList data={data} select={this.selectChannel}/>
					<RaidList data={data}/>
					<MemberList data={data}/>	
				</div>
			</div>
		);
	},
});

React.render(<App />, document.getElementById('app'));
