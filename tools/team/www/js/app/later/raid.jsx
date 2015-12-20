React = require('react/addons');
Dispatcher = require('../lib/dispatcher.jsx');
Routing = require('aviator');

module.exports = React.createClass({
	click: function(e) {
		Routing.navigate(
			"/:section/:channel/:raid", 
			{ namedParams: {
				section: this.props.routing.params.a,
				channel: this.props.routing.params.b,
				raid: this.props.data.uuid,
			} }
		);
		e.preventDefault();
	},
	render: function() {
		className = "raid";
		if ( this.props.selected == this.props.data.uuid ) {
			className = className + " active";
		}
		var now = new Date();
		var then = new Date(this.props.data.raid_time);
		

		var display = 'inline-block';
		this.raidTitle = this.props.data.raid_title;
		if (this.raidTitle == "") {
			this.raidTitle = this.props.data.name;
			this.dateString = "";
			display = 'none';
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
			} else if (hours == 12) {
				ampm = "pm"
			}
			var minutes = then.getMinutes();
			if (minutes < 10) {
				minutes = "0" + minutes;
			}

			this.dateString = month + "/" + date + " " + hours + ":" + minutes + ampm;
		}

		var dateLabelStyle = {
			width: '8.095238095238095em',
			display: display,
			textAlign: 'left'
		}

		return (
			<div className={className}>
				<div className="row">
					<div className="col-xs-12">
						<a onClick={this.click} className="btn btn-small btn-default btn-block pull-left" href="#">
							<div className="row">
								<div className="timeTitle col-xs-11">
									<span className="raidTime label label-primary" style={dateLabelStyle}>{this.dateString}</span>&nbsp;
									<span className="raidTitle">{this.raidTitle}&nbsp;</span>
								</div>
								<div className="members col-xs-1">
									<span className={'badge'}>{this.props.number}</span>	
								</div>
							</div>	
						</a>
						
					</div>
					
				</div>
			</div>);
	}
});
