React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');

module.exports = React.createClass({
	click: function(e) {
		Dispatcher.dispatch({actionType: "set", key: "raid", value: this.props.data.uuid});
		e.preventDefault();
	},
	render: function() {
		className = "raid";
		if ( this.props.selected == this.props.data.uuid ) {
			className = className + " active";
		}
		var now = new Date();
		var then = new Date(this.props.data.raid_time);
		
		/* this code no longer used, but leaving it here just in case
		var agoString = "";
		var seconds = Math.round( (then.getTime()/1000) - (now.getTime()/1000) );
		
		if (seconds < 0) {
			agoString = "";
		} else if (seconds < 60) {
			agoString = seconds + "s ";
		} else {
			var minutes = (seconds / 60).toFixed(1);

			if (minutes < 60) {
				agoString = minutes + "m ";
			} else {
				var hours = (minutes / 60).toFixed(1);

				if (hours < 24) {
					agoString = hours + "h ";
				} else {
					var days = (hours / 24).toFixed(1);
					agoString = days + "d ";
				}
			}
		}
		*/

		this.raidTitle = this.props.data.raid_title;
		if (this.raidTitle == "") {
			this.raidTitle = this.props.data.name;
			this.dateString = "";
		} else {
			var month = then.getMonth()*1 +1;
			var date = then.getDate();
			var ampm = "am";
			var hours = then.getHours();
			if (hours == 0) {
				hours = 12;
				ampm = "am";
			} else if (hours > 12) {
				hours = (hours - 12);
				ampm = "pm";	
			}
			var minutes = then.getMinutes();
			if (minutes < 10) {
				minutes = "0" + minutes;
			}

			this.dateString = month + "/" + date + " " + hours + ":" + minutes + ampm;
		}

		return (
			<div className={className}>
				<div className="row">
					<div className="col-md-12">
						<a onClick={this.click} className="btn btn-small btn-default btn-block pull-left" href="#">
							<span className="pull-left">
								<span className="label label-primary">{this.dateString}</span>&nbsp;
								{this.raidTitle} &nbsp;
								<span className={'badge'}>{this.props.number}</span>
							</span>
						</a>
						
					</div>
					
				</div>
			</div>);
	}
});
