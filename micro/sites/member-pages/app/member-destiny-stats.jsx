var React = require('react');
var dispatcher = require('./dispatcher.js');

var StatRow = React.createClass({
	click: function(e) {
		e.preventDefault();
		dispatcher.dispatch({
			type: "to",
			route: "lboard",
			params: {
				stat: this.props.name,
				type: "pvp",
			}
		})
	},
	render: function() {
		
		for ( var memberName in this.props.stat ) {
			if ( memberName != this.props.member.gamertag ) {
				continue;
			}

			var place = 1;
			var total = 0;
			var number = this.props.stat[memberName];

			for ( var otherMemberName in this.props.stat ) {
				total = total + 1;
				if ( this.props.stat[otherMemberName] > number ) {
					place = place + 1;
				}
			}

			var pct = Math.round( ( place / total ) * 100 );
			var widthl = (100 - pct) + "%";
			var widthr = pct + "%";
			return(
				<div key={this.props.name} className="row">
					<div className="col-md-3 col-md-offset-1 stat header">
						<a href="#" onClick={this.click}>
							{this.props.name}:
						</a>
					</div>
					<div className="col-md-2">
						{this.props.stat[memberName]}
					</div>
					<div className="col-md-1">
						{place}
					</div>
					<div className="col-md-2">
						<div style={{
							zIndex: "-1",
							position: "absolute",
							left: "0px",
							top: "5%",
							width: widthl,
							height: "90%",
							background: "#ddd",
						}}/>
						<div style={{
							zIndex: "-1",
							position: "absolute",
							left: widthl,
							top: "5%",
							width: widthr,
							height: "90%",
							background: "#eee",
						}}/>
						{pct}%
					</div>
				</div>
			);

		}
	},
});

var StatHeader = React.createClass({
	render: function() {
		return (
			<div key="header" className="row">
				<div className="col-md-3 col-md-offset-1 stat header">
					Stat Name
				</div>
				<div className="col-md-2 stat header">
					Stat Value
				</div>
				<div className="col-md-1 stat header">
					Rank
				</div>
				<div className="col-md-2 stat header">
					Percentile
				</div>
			</div>
		);
	}
});

module.exports = React.createClass({
	render: function() {
		var mainStats = [ (<StatHeader/>) ];

		for ( var statName in this.props.state["lb.pvp"] ) {
			mainStats.push((
				<StatRow
					key={statName}
					name={statName}
					stat={this.props.state["lb.pvp"][statName]}
					member={this.props.member}/>
			));
		}

		return (<div>{mainStats}</div>);
	},
});

