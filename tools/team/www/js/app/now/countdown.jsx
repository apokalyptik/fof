React = require('react/addons');

var CountdownTimer = React.createClass({
	getInitialState: function() {
		return {
			secondsInitial: 0,
			secondsRemaining: 0,
			ticked: false,
		};
	},
	tick: function() {
		this.setState({secondsRemaining: this.state.secondsRemaining - 1, ticked: true});
		if (this.state.secondsRemaining <= 0) {
			Dispatcher.dispatch({
				actionType: "lfg-looking",
				value: false
			});
			clearInterval(this.interval);
		}
	},
	componentDidMount: function() {
		this.setState({
			secondsInitial: this.props.secondsRemaining,
			secondsRemaining: this.props.secondsRemaining
		});
		this.interval = setInterval(this.tick, 1000);
	},
	componentWillUnmount: function() {
		clearInterval(this.interval);
	},
	render: function() {
		var niceDisplay = "";
		var seconds = this.state.secondsRemaining;
		var minutes = Math.floor( seconds / 60 );
		seconds = seconds - ( 60 * minutes ); 
		if ( this.state.secondsRemaining >= 3600 ) {
			var hours = Math.floor(this.state.secondsRemaining/3600)
			minutes = minutes - ( hours * 60 );
			if ( minutes == 60 ) {
				minutes = 0;
				hours = hours + 1;
			}
			niceDisplay = hours + "h";
		}
		niceDisplay = niceDisplay + minutes + "m" + seconds + "s"
		var pct = 0;
		if ( this.state.ticked ) {
			if ( this.state.secondsRemaining == this.state.secondsInitial ) {
				pct = 100;
			} else {
				if ( this.state.secondsInitial != 0 ) {
					pct = 100 - Math.ceil(
						100 * ( this.state.secondsRemaining / this.state.secondsInitial )
					);
				}
			}
		}
		return (
			<div className="center">
				<div style={{textShadow: "0 0 1px #fff"}}>{niceDisplay}</div>
				<div className="progress-bar" style={{textAlign:"left", marginTop:"-1.7em"}}>
					<span className="center" style={{width: pct+"%"}}></span>
				</div>
			</div>
		);
	}
});

module.exports = CountdownTimer;
