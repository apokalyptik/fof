var React = require('react');

module.exports = React.createClass({
	render: function() {

		var mainStats = [];

		for ( var statName in this.props.state["lb.pvp"] ) {

			for ( var memberName in this.props.state["lb.pvp"][statName] ) {

				if ( memberName != this.props.member.gamertag ) {
					continue;
				}

				var place = 1;
				var total = 0;
				var number = this.props.state["lb.pvp"][statName][memberName];

				for ( var otherMemberName in this.props.state["lb.pvp"][statName] ) {
					total = total + 1;
					if ( this.props.state["lb.pvp"][statName][otherMemberName] > number ) {
						place = place + 1;
					}
				}

				var pct = Math.round( ( place / total ) * 100 );

				mainStats.push((
					<div key={statName} className="row">
						<div className="col-md-5 col-md-offset-1">
							<strong>{statName}</strong>:
						</div>
						<div className="col-md-2">
							{this.props.state["lb.pvp"][statName][memberName]}
						</div>
						<div className="col-md-2 col-md-offset-1">
							top {pct}%
						</div>
					</div>
				));

			}

		}

		return (<div>{mainStats}</div>);

	},
});

