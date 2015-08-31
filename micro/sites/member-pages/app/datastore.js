module.exports = {
	data: {
		route: { name: "", data: {} },
		loaded: 0,
	},
	callbacks: [],
	listen: function(callback) {
		this.callbacks.push(callback)
	},
	notify: function() {
		for ( var i=0; i<this.callbacks.length; i++ ) {
			this.callbacks[i](this.data);
		}
	},
	set: function(k, v) {
		this.data[k] = v;
		if ( k == "loaded" ) {
			this.save();
		}
		this.notify();
	},
	save: function() {
		if ( window.localStorage ) {
			localStorage.setItem("cacheData", JSON.stringify(this.data));
		}
	}
};

if ( window.localStorage ) {
	var raw = localStorage.getItem("cacheData");
	if ( raw ) {
		var data = JSON.parse(raw);
		if ( data ) {
			module.exports.data = data;
		}
	}

}
