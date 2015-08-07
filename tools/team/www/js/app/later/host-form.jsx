React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
DateTimePicker = require('react-widgets/lib/DateTimePicker');

module.exports = React.createClass({
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
		if ( this.state.raidName == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "please enter an event name"});
			return;
		}

		if ( this.state.raidDateTimeString == null ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "please select a date and time"});
			return;
		} 

		jQuery.post("/rest/raid/host", {
			channel: 	this.state.channel,
			raid: 		this.state.raid,
			raidName: 	this.state.raidName,
			time: 		this.state.raidTimeSeconds,
			dateString: this.state.raidDateTimeString
		})
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
		this.setState({ "raidName": event.target.value })
		this.setState({ "raid": "[" + this.state.raidDateTimeString + "] " + this.state.raidName});
	},
	handleChannel: function(event) {
		this.setState({ "channel": event.target.value })
	},
	handleDateTime: function(value) {

		var date = new Date(value*1);

		// client based timezone. Good idea, bad idea?
		var timeZone = date.toString().match(/\(([A-Za-z\s].*)\)/)[1];
		var month = (date.getMonth() +1);
		var day = date.getDate()*1;
		var ampm = "am";
		var hours = date.getHours()*1;
		if (hours == 0) {
			hours = 12;
		} else if (hours > 12) {
			hours = hours - 12;
		}

		minutes = date.getMinutes()*1;

		if (minutes < 10) { 
			minutes = "0" + minutes;
		}

		var dateString = month + "/" + day + " " 
			+ hours + ":" +  minutes + ampm + " " + timeZone;

		this.setState({"raidDateTimeString": dateString});
		this.setState({"raidTimeSeconds": value});
		this.setState({"raid": "[" + this.state.raidDateTimeString + "] " + this.state.raidName});

	},
	handleDateTimeClick: function(event){
		if ($(event.target).hasClass("rw-input")){
			$("button.rw-btn-calendar").click();
		}
	},
	componentDidMount: function(){
		//make date/time input field readOnly
		$("input[type=text].rw-input").prop("readonly",true);
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

		var currentDate = new Date();
		var sevenDays = new Date(currentDate.getTime() + (7 * 24 * 60 * 60 * 1000) )
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
					<label htmlFor="DateTimePicker">Date and Time:</label>
					<DateTimePicker onChange={this.handleDateTime} onClick={this.handleDateTimeClick} min={currentDate} max={sevenDays}/>
					<small>Click the calendar to select a date, and the clock to select a time. Timezone will be selected by your browser.</small>
				</div>
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

