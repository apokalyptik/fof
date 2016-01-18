require('../../css/style.scss'); // So that webpack finds the scss file, and compiles it...

var Dispatcher = require( './lib/dispatcher.jsx');
var React = require('react/addons');
var Datastore = require( './lib/datastore.jsx');
var Hello = require('./hello/hello-main.jsx');
var SelectAnApp = require('./hello/select-an-app.jsx');
var LFGSelectGame = require('./now/select-game.jsx');
var LFGApp = require('./now/now-main.jsx');
var TeamApp = require('./later/later-main.jsx');
var MyLater = require('./later/my-raids.jsx');
var Notification = require("react-notification");
var Report = require('./report-a-friend/main.jsx');
var Routing = require('aviator');

Routing.pushStateEnabled = false;
Routing.setRoutes({
	target: {
		onChange: function( req, opts ) {
			var defaults = { a:"", b:"", c:"" };
			for ( var i in defaults ) {
				if ( "undefined" === typeof req.params[i] ) {
					req.params[i] = defaults[i]
				}
			}
			Datastore.set({routing: { params: req.params, opts: opts } });
		},
	},
	"/:a": {
		"/:b": {
			"/:c": {}
		},
	},
	"/*": "onChange",
});
Routing.dispatch();

var App = React.createClass({
	getInitialState: function() {
		return Datastore.data
	},

	componentDidMount: function() {
		jQuery.getJSON("/rest/login/check")
			.done(function(data) {
				Dispatcher.dispatch({actionType: "set", key: "cmd", value: data.cmd});
				if ( typeof data.username == "string" && data.username != "" ) {
					Dispatcher.dispatch({actionType: "set", key: "cmd", value: data.cmd});
					Dispatcher.dispatch({actionType: "set", key: "username", value: data.username});
					Dispatcher.dispatch({actionType: "username", value: data.username});
					Dispatcher.dispatch({actionType: "set", key: "authenticated", value: true});
				}
				Datastore.subscribe(this.acceptData)
				Dispatcher.dispatch({actionType: "set", key: "checked", value: true});
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
    getErrorNotificationStyles: function() {

        return  {
        	bar: {
        		backgroundColor: '#ff9999'
	        }, 
	        active: {
	        	left: '3rem'
	        }, 
	        action: {
	        	color: '#ff9999'
	        }
	    };
    },
    handleNotificationClick: function(notification) {
    	Dispatcher.dispatch({actionType: "set", key: notification, value: ""});
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
							Please use the slack command &ldquo;<strong>/team</strong>&rdquo; to log in
							</h2>
						</div>
					</div>
				</div>
			);
		}

		var crumbs = [
			( <div key="appselect" className="box col-xs-2 nopadding">
				  <SelectAnApp key="selectanapp" routing={this.state.routing}/>
			  </div> )
		];

		var WorkSpace;
		switch ( this.state.routing.params.a ) {
			case "later":
				crumbs.push((<MyLater key="mylater" state={this.state} routing={this.state.routing}/>));
				WorkSpace = ( <TeamApp state={this.state} routing={this.state.routing}/> );
				break;
			case "now":
				WorkSpace = ( <LFGApp state={this.state.lfg} routing={this.state.routing}/> );
				break;
			case "report":
				WorkSpace = ( <Report state={this.state.report} routing={this.state.routing}/> );
				break;
			default:
				return ( <Hello routing={this.state.routing}/> );
		}

		var Error;
		if ( this.state.error ) {
			Error = (
				<Notification 
					ref="errorNotification"
					isActive={true} 
					message={this.state.error} 
					action="&times;" 
					style={{
			            bar: {
			              top: '1rem',
                          bottom: 'auto',
			              font: '1.25rem normal Roboto, sans-serif',
			              backgroundColor: '#CC0000',
			              color: '#FFFFFF',
			              zIndex: 9999
			            },
			            action: {
			              color: '#FFFFFF',
                          fontSize: '1.25rem'
			            }
			        }} 
					dismissAfter={30000}
					onDismiss={this.handleNotificationClick.bind(null,'error')}
					onClick={this.handleNotificationClick.bind(null,'error')}/>
			);
		}

		var Success;
		if ( this.state.success ) {
			Success = (
				<Notification
					isActive={true}
					message={this.state.success}
					style={{
			            bar: {
			              top: '1rem',
			              bottom: 'auto',
			              font: '1.25rem normal Roboto, sans-serif',
			              backgroundColor: '#ADEBAD',
			              color: '#2C6710',
			              border: 'solid 1px #9CBF9C',
			              zIndex: 9999
			            },
			            action: {
			              color: 'rgb(0, 0, 0)',
			              fontSize: '1.25rem'
			            }
			        }} 
			        action="&times;"
					dismissAfter={10000}
					onDismiss={this.handleNotificationClick.bind(null,'success')}
					onClick={this.handleNotificationClick.bind(null,'success')}/>
			);
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
								<div className="breadcrumbs-lgr row nomargin">
									{crumbs}
									<div className="rt col-xs-1 nopadding"/>
									<button className="hidden-xs btn btn-default pull-right" style={{
										border: "0.15em solid #1b1c52", 
										marginTop: "0.1em",
										marginRight: "0.4em"
									}} onClick={function(e){
										e.preventDefault();
										Routing.navigate("/report");
									}}>Report a Claim</button>
								</div>
							</div>
							<div className="notices">
								{Error}
								{Success}
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
