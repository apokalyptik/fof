var Immutable = require('immutable');

var state = {
	data: {},
	cbs: [],
	register: function(cb) {
		this.cbs.push(cb)
	},
	flush: function() {
		var data = this.get();
		for( var i=0; i<this.cbs.length; i++ ) {
			this.cbs[i](data)
		}
	},
	set: function(o) {
		for( var i in o ) {
			this.data[i] = o[i];
		}
	},
	get: function() {
		return Immutable.Map(this.data)
	}
}

var Dispatcher = require('flux').Dispatcher;
var stateDispatcher = new Dispatcher();
stateDispatcher.register(function(payload) {
	switch( payload.action ) {
		case "register":
			state.register(payload.callback);
			return;
		case "set":
			state.set(payload.val);
			state.flush();
			return;
	}
});

module.exports = {
	get:      function()   { return state.get()                                           },
	set:      function(o)  { stateDispatcher.dispatch({action: "set", val: o})            },
	register: function(cb) { stateDispatcher.dispatch({action: "register", callback: cb}) }
};
