var active = { authenticated: false, channel: null, raid: null };

var returnFalse = function( e ) { e.stopPropagation(); return false; }

var Channel = React.createClass({displayName: "Channel",
	render: function() {
		var classes = [ "raidChannel" ];
		if ( this.props.name == active.channel ) {
			classes.push("active");
		}
		return(
			React.createElement("div", {className: classes.join(" ")}, 
				React.createElement("a", {href: "#", onClick: this.props.select(this.props.name)}, this.props.name)
			)
		);
	},
});

var ChannelList = React.createClass({displayName: "ChannelList",
	render: function() {
		var channelList = (
			React.createElement("strong", null, "There are currently no raids being hosted")
		);
		var channels = [];
		for ( var c in this.props.data ) {
			if ( c[0] != "@" ) {
				channels.push(c)
			}
		}
		if ( channels.length > 0 ) {
			var channelList = channels.map(function (channelName) {
				return ( React.createElement(Channel, {name: channelName, select: this.props.select}) );
			});
		}
		return(
			React.createElement("div", {className: "col-md-3"}, 
				React.createElement("h2", null, "Channels"), 
				channelList
			)
		);
	},
});

var RaidList = React.createClass({displayName: "RaidList",
	render: function() {
		var raidList = (
			React.createElement("strong", null, 
				"Please select a channel to see the raid list"
			)
		);
		return(
			React.createElement("div", {className: "col-md-5"}, 
				React.createElement("h2", null, "Raids"), 
				raidList
			)
		);
	},
});

var MemberList = React.createClass({displayName: "MemberList",
	render: function() {
		var memberList = (
			React.createElement("strong", null, 
				"Please select a raid to see the member list and be able to join or part"
			)
		);
		return(
			React.createElement("div", {className: "col-md-4"}, 
				React.createElement("h2", null, "Members"), 
				memberList
			)
		);
	},
});

var LoginInit = React.createClass({displayName: "LoginInit",
	render: function() {
		return(
			React.createElement("div", {className: "col-md-3 col-md-offset-4"}, 
				React.createElement("form", null, 
					React.createElement("div", {className: "form-group"}, 
					React.createElement("h2", null, "Please Log In")
					), 
					React.createElement("div", {className: "form-group"}, 
						React.createElement("label", {htmlFor: "inputLogin"}, 
							"Slack username or email address"
						), 
						React.createElement("input", {id: "inputLogin", type: "text", value: "", className: "form-control"})
					), 
					React.createElement("button", {type: "submit", className: "btn btn-default"}, "Submit")
				)
			)
		);
	}
});

var Login = React.createClass({displayName: "Login",
	getInitialState: function() {
		return {
			username: "",
			authenticated: false,
			step: 0,
		};
	},
	render: function() {
		if ( this.state.step == 0 ) {
			return (React.createElement(LoginInit, null));
		}
	},
});

var App = React.createClass({displayName: "App",
	getIntialState: function() {
		return {};
	},
	selectChannel: function(name) {
		active.channel = name;
		this.setState(this.state)
	},
	render: function() {
		if ( active.authenticated == false ) {
			return( 
				React.createElement("div", {className: "container-fluid"}, 
					React.createElement("div", {className: "row"}, 
						React.createElement(Login, null)
					)
				)
			);
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
			React.createElement("div", {className: "container-fluid"}, 
				React.createElement("div", {className: "row"}, 
					React.createElement(ChannelList, {data: data, select: this.selectChannel}), 
					React.createElement(RaidList, {data: data}), 
					React.createElement(MemberList, {data: data})	
				)
			)
		);
	},
});

React.render(React.createElement(App, null), document.getElementById('app'));
