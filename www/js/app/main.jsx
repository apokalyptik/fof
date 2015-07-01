var Dispatcher = require( './lib/dispatcher.jsx');
var React = require('react/addons');
var Datastore = require( './lib/datastore.jsx');
var Hello = require('./hello/hello-main.jsx');
var SelectAnApp = require('./hello/select-an-app.jsx');
var LFGSelectGame = require('./now/select-game.jsx');
var LFGApp = require('./now/now-main.jsx');
var TeamApp = require('./later/later-main.jsx');
var MyLater = require('./later/my-raids.jsx');

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
		if ( this.state.checked == false ) {
			return (<div/>);
		}

		if ( this.state.authenticated == false ) {
			return(
				<div className="container-fluid">
					<div className="row">
						<div className="col-md-6 col-md-offset-3 center">
							<h2 className="dark">
							please use the slack command &ldquo;<strong>/team</strong>&rdquo; to log in
							</h2>
						</div>
					</div>
				</div>
			);
		}

		if ( this.state.viewing == "hello" ) {
			return ( <Hello/> );
		}

		var crumbs = [
			( <li key="appselect" className="box">
				  <SelectAnApp key="selectanapp" viewing={this.state.viewing}/>
			  </li> )
		];

		crumbs.push((<MyLater key="mylater" state={this.state}/>));

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

		var Error;
		if ( this.state.error ) {
			Error = (
				<DisplayFor seconds={30} message={(<div className="error">{this.state.error}</div>)}/>
			);
		}

		var Success;
		if ( this.state.success ) {
			Success = (
				<DisplayFor seconds={10} message={(<div className="success">{this.state.success}</div>)}/>
			);
		}

		return(
			<div>
				<div className="notices">
					{Error}
					{Success}
				</div>
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

Dispatcher.register(function(payload) {
	var doReRender = false;
	switch ( payload.actionType ) {
		case "serverStateUpdate":
			for ( var i in payload.data ) {
				switch( i ) {
					case "lfg":
						Datastore.data.lfg.username = Datastore.data.username;
						Datastore.data.lfg.prevlfg = Datastore.data.lfg.lfg;
						Datastore.data.lfg.lfg = payload.data[i];
						break;
					default:
						Datastore.data[i] = payload.data[i];
						break;
				}
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
			break;
		case "mset":
			Datastore.setThings(payload.what);
			break;
		case "set":
			Datastore.setThing(payload.key, payload.value);
			break;
	}
});

jQuery(document).ready(function() {
	React.render(<App />, document.getElementById('app'));
})
