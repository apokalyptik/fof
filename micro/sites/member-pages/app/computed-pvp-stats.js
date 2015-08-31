module.exports = function(data) {
	var newStats = {};
	var newStatname = "";
	var otherStat = "";
	for ( var statName in data ) {
		switch ( statName ) {
			case "kills":
				newStatname = "killDeathRatio";
				otherStat = "deaths";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				
				newStatname = "killsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;

				}
				
				newStatname = "killsPerMinute";
				otherStat = "secondsPlayed";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / (b /60))*100)/100;
				}
				
				break;
			case "activitiesWon":
				// ( activitiesWon / activitiesEntered ) * 100
				newStatname = "winRatio";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;

				}
				break;
			case "deaths":
				// deaths / activitiesEntered
				newStatname = "deathsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				
				// deaths / secondsPlayed/60
				newStatname = "deathsPerMinute";
				otherStat = "secondsPlayed";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / (b/60))*100)/100;
				}
				break;
			case "orbsDropped":
				// orbsDropped / activitiesEntered
				newStatname = "orbsMadePerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "orbsGathered":
				// orbsGathered / activitiesEntered
				newStatname = "orbsGrabbedPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "precisionKills":
				// precisionKills / activitiesEntered
				newStatname = "PrecisionKillsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "weaponKillsGrenade":
				// weaponKillsGrenade / activitiesEntered
				newStatname = "grenadeKillsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "weaponKillsMelee":
				// weaponKillsMelee / activitiesEntered
				newStatname = "meleeKillsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "weaponKillsSuper":
				// weaponKillsSuper / activitiesEntered
				newStatname = "superKillsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "zonesCaptured":
				// zonesCaptured / activitiesEntered
				newStatname = "zoneCapsPerActivity";
				otherStat = "activitiesEntered";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / b)*100)/100;
				}
				break;
			case "assists":
				// assists / secondsPlayed/60
				newStatname = "assistsPerMinute";
				otherStat = "secondsPlayed";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / (b/60))*100)/100;
				}
				break;
			case "weaponKillsSuper":
				// secondsPlayed/60 / weaponKillsSuper
				newStatname = "supersPerMinute";
				otherStat = "secondsPlayed";
				newStats[newStatname] = {};
				for ( var memberName in data[statName] ) {
					var a = data[statName][memberName];
					var b = data[otherStat][memberName];
					newStats[newStatname][memberName] = Math.round((a / (b/60))*100)/100;
				}
				break;
		}
	}
	for ( var i in newStats ) {
		data[i] = newStats[i];
	}
	return data;
};
