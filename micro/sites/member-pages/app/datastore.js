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
		this.notify();
	}
};

