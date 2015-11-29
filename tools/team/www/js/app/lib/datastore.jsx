module.exports = {
	callbacks: [],
	data: {
		authenticated: false,
		checkedUsername: false,
		username: "",
		checked: false,
		command: "",
		updated_at: "",
		hosting: false,
		channels: [],
		lfg: {
			username: "",
			my: {},
			prevlfg: {},
			lfg: {},
			time: 120,
			looking: false,
		},
		error: null,
		success: null,
	},
	set: function(obj) {
		for( var i in obj ) {
			this.data[i] = obj[i];
		}
		this.emitChange();
	},
	setThings: function(what) {
		for( var i=0; i<what.length; i++ ) {
			this.data[what[i].key] = what[i].value;
		}
		this.emitChange();
	},
	setThing: function(thing, value) {
		this.data[thing] = value;
		this.emitChange();
	},
	subscribe: function(callback) {
		this.callbacks.push(callback);
	},
	emitChange: function() {
		for( var i = 0; i < this.callbacks.length; i++ ) {
			this.callbacks[i]( this.data );
		}
	}
}

