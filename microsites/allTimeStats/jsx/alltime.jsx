var React = require("react/addons")
var FluxDispatcher = require('flux').Dispatcher;
var Dispatcher = new FluxDispatcher();

var Store = {
	data: {
		possible: null,
		section: "merged",
		stat: "",
		statData: null,
		nameFilter: "",
	},
	callbacks: [],
	subscribe: function(callback) {
		this.callbacks.push(callback)
	},
	notify: function() {
		for ( var i=0; i<this.callbacks.length; i++ ) {
			this.callbacks[i](this.data);
		}
	},
	set: function(k, v) {
		this.data[k] = v;
		this.notify();
	}
};

Dispatcher.register(function(e) {
	switch( e.type ) {
		case "nameFilter":
			Store.set("nameFilter", e.value);
			break;
		case "setSection":
			Store.data.stat = "";
			Store.data.statData = null;
			Store.set("section", e.value);
			break;
		case "setStat":
			if ( e.value == "" ) {
				Store.data.statData = null;
				Store.set("stat", "");
				return
			}
			Store.data.stat = e.value;
			$.getJSON("http://fofgaming.com:8880/destiny/stats/alltime/"+Store.data.section+"/"+Store.data.stat+".json")
				.done(function(data){
					Store.set("statData", data);
				})
				.fail(function() {
					alert("Sorry, something went wrong trying to fetch this data...");
					Store.notify();
				});
			break;
		case "setPossible":
			Store.set("possible", e.value);
			break;
	}
});

var PossibleStats = React.createClass({
	handleChange: function(e) {
		Dispatcher.dispatch({ type: "setStat", value: e.target.value })
	},
	shouldComponentUpdate: function() {
		return true;
	},
	render: function() {
		var options = [ (<option name="none" key="none" value="">Select a Stat</option>) ];
		if ( typeof this.props.possible[this.props.section] != "undefined" ) {
			for( var i=0; i<this.props.possible[this.props.section].length; i++ ) {
				var stat = this.props.possible[this.props.section][i];
				options.push( (<option name={stat} key={stat} value={stat}>{stat}</option>) );
			}
		}
		return(
			<select className="form-control" value={this.props.selected} onChange={this.handleChange}>{options}</select>
		);
	},
});

var PossibleSections = React.createClass({
	handleChange: function(e) {
		Dispatcher.dispatch({ type: "setSection", value: e.target.value })
	},
	render: function() {
		var options = [];
		for ( var k in this.props.possible ) {
			options.push((<option name={k} key={k} value={k}>{k}</option>));
		}
		return(<select className="form-control" onChange={this.handleChange} value={this.props.selected}>{options}</select>);
	},
});

var Stats = React.createClass({
	render: function() {
		var rows = [(
			<tr key="header">
				<td><strong>Rank</strong></td>
				<td><strong>Member Name</strong></td>
				<td><strong>{this.props.stat}</strong></td>
			</tr>
		)];
		var re = new RegExp("(" + this.props.nameFilter + ")", "i");
		for( var i=0; i<this.props.data.length; i++ ) {
			if ( !this.props.data[i].member.match(re) ) {
				continue;
			}
			rows.push((
				<tr key={this.props.data[i].member}>
					<td>{i + 1}</td>
					<td>{this.props.data[i].member}</td>
					<td>{this.props.data[i].value}</td>
				</tr>
			));
		}
		return(<div><hr/><table className="table table-striped"><tbody>{rows}</tbody></table></div>);
	},
});

var App = React.createClass({
	changeNameFilter: function(e) {
		Dispatcher.dispatch({type: "nameFilter", value: e.target.value})
	},
	render: function() {
		if ( this.state.possible == null ) {
			return(<div/>);
		}
		var stats = null;
		if ( this.state.statData != null ) {
			stats = (
					<Stats
						key="stats"
						nameFilter={this.state.nameFilter}
						data={this.state.statData}
						stat={this.state.stat}/>
					);
		}
		return (
			<div>
				<div className="row"><div className="col-md-12">
					<h2>Select a stat type and specific stat to get the goods</h2>
					<p>You can also filter the names by typing into the box</p>
				</div></div>
				<div className="row">
					<div className="col-md-2" style={{textAlign:"right"}}>
						<PossibleSections selected={this.state.section} possible={this.state.possible}/>
					</div>
					<div className="col-md-3">
						<PossibleStats selected={this.state.stat} possible={this.state.possible} section={this.state.section}/>
					</div>
					<div className="col-md-3">
						<input
							className="form-control"
							value={this.state.filterNames}
							onChange={this.changeNameFilter}
							placeholder="Type a name"/>
					</div>
				</div>
				<div className="row">
					<div className="col-md-12">
						{stats}
					</div>
				</div>
			</div>
		);
	},
	getInitialState: function() {
		return Store.data;
	},
	componentDidMount() {
		Store.subscribe(this.setState.bind(this))
		$.getJSON('http://fofgaming.com:8880/destiny/stats/alltime/keys.json')
			.done(function(data) {
				Dispatcher.dispatch({ type: "setPossible", value: data })
			})
			.fail(function() {
				alert("Failed to fetch enough data to render this page, sorry");
			});
	},
});

$(document).ready(function() {
	React.render(
		React.createElement(App, null),
		document.getElementById('container')
	);
});
