Dispatcher = require('../lib/dispatcher.jsx');
React = require('react');

module.exports = React.createClass({
	getInitialState: function() {
		return { about: "", message: "", enabled: true, members: [] };
	},
	componentWillMount: function() {
		jQuery.getJSON("http://team.fofgaming.com:8880/fof/members.json")
			.done(function(data) {
				this.setState({members: data});
			}.bind(this))
			.fail(function() {
				window.setTimeout(this.componentWillMount, 1000);
			}.bind(this));
	},
	submit: function() {
		if ( this.state.about == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please choose a member"});
			return;
		}
		if ( this.state.message == "" ) {
			Dispatcher.dispatch({actionType: "set", key: "error", value: "Please supply a message"});
			return;
		}
		Dispatcher.dispatch({actionType: "set", key: "error", value: ""});
		Dispatcher.dispatch({actionType: "set", key: "success", value: "Submitting your report now..."});
		this.setState({enabled: false});
		jQuery.post("/rest/report", {
			about:   this.state.about,
			message: this.state.message
		})
		.done(function() {
			Dispatcher.dispatch({actionType: "set", key: "success", value: "Report submitted!"});
			this.setState({ about: "", message: "", enabled: true });
		}.bind(this))
		.fail(function() {
			Dispatcher.dispatch({actionType: "set", key: "success", value: ""});
			Dispatcher.dispatch({
				actionType: "set", 
				key: "error", 
				value: "There was an error sending your report. Please try again in a few minutes"
			});
			this.setState({enabled: true});
		}.bind(this))
	},
	render: function() {
		var disabled = true;
		if ( this.state.enabled === true )  {
			disabled = false;
		}
		var members = [
			(<option value="" key="-1">Who is behaving badly?</option>),
		];
		for ( var i=0; i<this.state.members.length; i++ ) {
			members.push((
				<option value={this.state.members[i].username} key={i}>{this.state.members[i].username}</option>
			));
		}
		return (
			<div className="container-fluid">
				<div className="row">
					<div className="col-md-6 col-md-offset-3">
						<h4>Report Bad Behavior</h4>
						<div className="form-group">
							<select className="form-control" id="about" name="about" value={this.state.about}
								placeholder="required"
								onChange={function(e) { this.setState({about: e.target.value}); }.bind(this)}>
									{members}
							</select>
						</div>
						
						<div className="form-group">
						<label htmlFor="message">What do you want to report about them?</label>
							<textarea className="form-control" name="message" id="message" rows="3"
								placeholder="required"
								onChange={function(e) { this.setState({message: e.target.value}); }.bind(this)}
								value={this.state.message}/>
						</div>

						<div className="form-group text-right">
							<input className="btn btn-default" type="submit" disabled={disabled} onClick={this.submit}/>
						</div>
						
					</div>
				</div>
			</div>
		);
	}
});
