var React = require('react/addons');
var Dispatcher = require('../lib/dispatcher.jsx');
var CountdownTimer = require('./countdown.jsx');

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
				alert("Pinged " +un);
				/*
				Dispatcher.dispatch({
					actionType: "set",
					key: "success",
					value: "Ping to " +un+ " successful"
				});
				Dispatcher.dispatch({
					actionType: "set",
					key: "failure",
					value: null
				});
				*/
			})
			.fail(function() {
				alert("Ping to " +un+ " unsuccessful");
				/*
				Dispatcher.dispatch({
					actionType: "set",
					key: "success",
					value: null,
				});
				Dispatcher.dispatch({
					actionType: "set",
					key: "failure",
					value: "Ping to " +un+ " unsuccessful"
				});
				*/
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
			var name = this.props.forWhat[i];
			if ( this.props.forWhat[name] == false ) {
				continue;
			}
			list.push(( this.renderSection(name) ) );
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

module.exports = LFGAppLooking;
