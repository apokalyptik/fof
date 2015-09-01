var React = require('react');
var dispatcher = require('./dispatcher.js');

var EntryRow = React.createClass({
	click: function(e) {
		e.preventDefault();
		dispatcher.dispatch({
			type: "to",
			route: "member",
			params: {
				name: this.props.name,
				game: "destiny",
				section: "pvp"
			}
		})
	},
	render: function() {
		return (
			<tr>
				<td><a onClick={this.click} href="#">{this.props.name}</a></td>
				<td>{this.props.value}</td>
			</tr>
		);
	},
});

function compare(a,b) {
	if (a.value > b.value)
		return -1;
	if (a.value < b.value)
		return 1;
	return 0;
}

module.exports = React.createClass({
	render: function() {
		var entries = [];
		var section = "lb." + this.props.state.route.data.type;
		var stats = this.props.state[section][this.props.state.route.data.stat];
		for ( var i in stats ) {
			entries.push( { name: i, value: stats[i] } );
		}
		entries.sort(compare);
		for ( var i=0; i<entries.length; i++ ) {
			entries[i] = (<EntryRow key={i} name={entries[i].name} value={entries[i].value}/>);
		}
		return (
			<div className="container-fluid">
				<div className="row">
					<div className="col-md-6 col-md-offset-3">
						<h3>{this.props.state.route.data.type} / {this.props.state.route.data.stat}</h3>
						<a href="#" onClick={function(e){e.preventDefault(); window.history.back();}}>Back</a>
					</div>
				</div>
				<div className="row">
					<div className="col-md-6 col-md-offset-3">
						<table className="table table-striped table-bordered">
							<thead>
								<tr>
									<td>User</td>
									<td>Value</td>
								</tr>
							</thead>
							{entries}
						</table>
					</div>
				</div>
			</div>
		);
	},
});
