React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
DateTimePicker = require('../lib/DateTimePicker.jsx');

module.exports = React.createClass({
	getInitialState: function() {
		Dispatcher.dispatch({actionType: "set", key: "error", value: ""});
		return {
			error: "",
			channel: "",
			raid: "",
		}
	},
	componentDidMount: function() {
		this.handleNewDateTimePicker();
	},
	submit: function(e) {
		if ( this.state.channel == null || this.state.channel == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please select a channel"});
			return;
		}

		if ( this.state.raidName == null) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please enter an event name"});
			return;
		}


		//build POST parameter values
		var raid = "[" + this.state.raidDateTimeString + "] " + this.state.raidName;


		jQuery.post("/rest/raid/host", {
			channel: 	this.state.channel,
			raid: 		raid,
			raidName: 	this.state.raidName,
			time: 		this.state.raidTime,
			timeZone:   this.state.timezoneString,
			dateString: this.state.raidDateTimeString
		})
			.done(function(data) {
				Dispatcher.dispatch({actionType: "set", key: "error", value: ""});
				Dispatcher.dispatch({actionType: "set", key: "success", value: "Your event has been created!"});
				this.props.cancel(e)
			}.bind(this))
			.fail(function(data) {
				if ( data.status == 403 ) {
					location.reload(true);
				}
				Dispatcher.dispatch({actionType: "set", key: "error", value: data.responseText});
			}.bind(this));
	},
	handleRaid: function(event) { 
		this.setState({ "raidName": event.target.value })
		// this.setState({ "raid": "[" + this.state.raidDateTimeString + "] " + this.state.raidName});
	},
	handleChannel: function(event) {
		this.setState({ "channel": event.target.value })
	},
	handleNewDateTimePicker: function(){
		var dateTimePicker = this.refs.dateTime;
		var eventDate = new Date(dateTimePicker.state.date);

		//set dateString
		var dateString = dateTimePicker.state.dateString;
		var hour = dateTimePicker.state.hourString*1;
		var minute = dateTimePicker.state.minuteString;
		var timeString = hour + ":" + minute + dateTimePicker.state.ampmString.toLowerCase();
		var timezoneString = dateTimePicker.state.timeZoneText;
		this.setState({
			"raidTime" : eventDate.getTime(),
			"dateString": dateString,
			"timeString": timeString,
			"timezoneString" : timezoneString,
			"raidDateTimeString": dateString.substring(0,dateString.lastIndexOf("/")) + " " + timeString + " " + timezoneString
		});
		
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
				<DateTimePicker ref="dateTime" maxDays="7" onChange={this.handleNewDateTimePicker}/>
				<div className="form-group">
					<label htmlFor="name">Name of your Event</label>
					
					<input 
						onChange={this.handleRaid}
						type="text" className="form-control" id="name" placeholder="Event Name"/>
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

