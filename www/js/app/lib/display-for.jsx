React = require('react/addons');

module.exports = React.createClass({
	getInitialState: function() {
		return {
			secondsRemaining: 0,
			message: "",
		};
	},
	tick: function() {
		this.setState({secondsRemaining: this.state.secondsRemaining - 1, ticked: true});
		if (this.state.secondsRemaining <= 0) {
			clearInterval(this.interval);
		}
	},
	componentDidMount: function() {
		this.setState({
			message: this.props.message,
			secondsRemaining: this.props.seconds
		});
		this.interval = setInterval(this.tick, 1000);
	},
	componentWillUnmount: function() {
		clearInterval(this.interval);
	},
	render: function() {
		if ( this.state.secondsRemaining < 1 ) {
			return (<span/>);
		}
		return (this.state.message);
	}
});

