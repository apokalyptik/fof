React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
DateTimePicker = require('react-widgets/lib/DateTimePicker');

module.exports = React.createClass({
	getInitialState: function() {
		Dispatcher.dispatch({actionType: "set", key: "error", value: ""});
		return {
			error: "",
			channel: "",
			raid: "",
		}
	},
	submit: function(e) {
		if ( this.state.channel == null || this.state.channel == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please select a channel"});
			return;
		}
		if ( this.state.dateString == null ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please select a date"});
			return;
		} 

		if ( this.state.timeString == null ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please select a time"});
			return;
		} 
		if ( this.state.raidName == null) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please enter an event name"});
			return;
		}

		var date = new Date(this.state.dateString + " " + this.state.timeString);

		var timeZone = this.getTimeZone(date);

		//build POST parameter values
		var dateString = this.state.dateString;
		dateString = dateString.substring(0,dateString.lastIndexOf("/"));
		var raidDateTimeString = dateString + ", " + this.state.timeString + " " + timeZone;
		var time =  date.getTime();
		var raid = "[" + raidDateTimeString + "] " + this.state.raidName;


		jQuery.post("/rest/raid/host", {
			channel: 	this.state.channel,
			raid: 		raid,
			raidName: 	this.state.raidName,
			time: 		time,
			timeZone:   timeZone,
			dateString: raidDateTimeString
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
				Dispatcher.dispatch({actionType: "set", key: "error", value: responseText});
			}.bind(this));
	},
	getTimeZone: function(date) {
		var offset = -( date.getTimezoneOffset() );
		var isDst =  offset < this.stdTimezoneOffset();
		var isUS = offset < 0;

		if (date instanceof Date) {
			
			var timezone = "";

			if (isUS) {
				// convert DST to ST
				if (isDst) {
					offset = offset - 1;
				}

				// US times
				switch(offset) {
					case -10:
						timezone = "H";
						break;
					case -9:
						timezone = "AK";
						break;
					case -8:
						timezone = "P";
						break;
					case -7:
						timezone = "M";
						break;
					case -6:
						timezone = "C";
						break;
					case -5:
						timezone = "E";
						break;
					default:
						return date.toString().match(/\(([A-Za-z\s].*)\)/)[1];
				}

				if (isDst) {
					return timezone + "DT";	// US Daylight Savings Time
				} else if (isUS) {
					return timezone + "ST";	// US Standard Time
				}
			} else {

				//europe
				switch (offset) {
					case 0:
						return "GMT";
					case 1:
						if (isDst) {
							return "BST";
						} else {
							return "CET";
						}
					case 2:
						if (isDst) {
							return "CEST";
						} else {
							return "EET";
						}
					case 3:
						if (isDst) {
							return "EEST";
						} else {
							return "MSK";
						}
					case 4:
						if (isDst) {
							return "MSD";
						} else {
							return "SAMT";
						}
					default:
						return date.toString().match(/\(([A-Za-z\s].*)\)/)[1];
				}

			}

		} else {
			return null;
		}
	},
	stdTimezoneOffset: function() {
		var now = new Date();
		var jan = new Date(now.getFullYear(), 0, 1);
		var jul = new Date(now.getFullYear(), 6, 1);
		return Math.max(jan.getTimezoneOffset(), jul.getTimezoneOffset());
	},
	handleRaid: function(event) { 
		this.setState({ "raidName": event.target.value })
		// this.setState({ "raid": "[" + this.state.raidDateTimeString + "] " + this.state.raidName});
	},
	handleChannel: function(event) {
		this.setState({ "channel": event.target.value })
	},
	handleDate: function(value) {

		var month = (value.getMonth() +1);
		var day = value.getDate()*1;
		var year = value.getFullYear()*1;
		var dateString = month + "/" + day + "/" + year;
		this.setState({"dateString": dateString});

		var timeString = this.state.timeString ==  null ? "12:00 am" : this.state.timeString;
		this.setState({"raidDateTimeString" : dateString + " " + timeString});

	},
	handleTime: function(value) {

		var ampm = "am";
		var hours = value.getHours()*1;
		if (hours == 0) {
			hours = 12;
		} else if (hours == 12) {
			ampm="pm";
		} else if (hours > 12) {
			hours = hours - 12;
			ampm="pm"
		}

		minutes = value.getMinutes()*1;

		if (minutes < 10) { 
			minutes = "0" + minutes;
		}
		var timeString = hours + ":" + minutes + " " + ampm;
		this.setState({"timeString":timeString});

		var dateString;
		if (this.state.dateString == null) {
			var now = new Date();
			var month = (now.getMonth() +1);
			var day = now.getDate()*1;
			var year = now.getFullYear()*1;
			dateString = month +"/"+ day + "/" + year;

		} else {
			dateString = this.state.dateString
		}
		this.setState({"raidDateTimeString" : dateString + " " + timeString});

	},
	handleDateClick: function(event){
		if ($(event.target).hasClass("rw-input")){
			$("#datePicker button.rw-btn-calendar").click();
		}
	},
	handleTimeClick: function(event){
		if ($(event.target).hasClass("rw-input")){
			$("#timePicker button.rw-btn-time").click();
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
				<div id="datePicker" className="form-group">
					<label htmlFor="DatePicker">Date:</label>
					<DateTimePicker time={false} format={"MMM dd, yyyy"} onChange={this.handleDate} onClick={this.handleDateClick} min={currentDate} max={sevenDays}/>
					<small>Click the calendar to select a date.</small>
				</div>
				<div id="timePicker" className="form-group">
					<label htmlFor="TimePicker">Time:</label>
					<DateTimePicker calendar={false} onChange={this.handleTime} onClick={this.handleTimeClick}/>
					<small>Click the clock to select a time. Timezone will be selected by your browser.</small>
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

