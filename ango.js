
angular.module('ango', [])
	.factory('Ango', function($q, $http) {
		var AngoService = {
			procedures: {},
			callbacks: {},
		};

		// Run is a fire-and-forget method to run a procedure with given data
		// This looks up a normal (no callback) procedure handler
		AngoService.Run = function(procedure, data) {
			//++
			return;
		}

		// Call sends a request and returns a promise.
		// The returned deferred is resolved by the first response from server.
		// Also allows server to send notifications (e.g. processing progress).
		// The given header
		AngoService.Call = function(procedure, data) {
			// create deferred
			var deferred = $q.defer();
			// register callback
			AngoService.callbackCounter++;
			cbid = AngoService.callbackCounter
			AngoService.callbacks[cbid] = deferred;

			var out = {
				type: "req",
				procedure: procedure,
				data: data,
				cb_id: cbid,
				def_id: cbid,
			};

			// all done
			return deferred.promise;
		};

		// RegisterProcedure allows you to implement client-side procedure that can be called by the server.
		// When the handler returns a promise and the server expects a callback the resolve/reject is sent to server
		// When the server expects a callback, but no promise is returns, an error is sent to the server.
		AngoService.RegisterProcedure = function(name, handler) {
			if(typeof name != "string") {
				throw("Invalid procedure name, not a string");
			}
			if(typeof handler != "function") {
				throw("Invalid procedure handler, not a function");
			}

			AngoService.procedures[name] = handler;
		}

		var ws = new WebSocket("ws://"+window.location.href.split("/")[2]+"/ango-websocket");

		ws.onopen = function(){
			console.log('ws opened');
		}

		ws.onmessage = function(message) {
			handler(JSON.parse(message.data));
		}

		ws.onerror = function() {
			console.log('ws error');
			//++ run hooks?
		}
		ws.onclose = function() {
			console.log('ws closed');
			//++ run hooks?
		}

		function sendPromiseFullfillment(type, def_id, data) {
			var out = {
				type: type,
				def_id: def_id,
				data: data,
			};
			ws.send(JSON.stringify(out));
		}

		function handler(msg) {
			switch(msg.type) {
				case "req":
					// lookup procedure
					var proc = AngoService.procedures[msg.procedure];
					if(typeof proc != "function") {
						// send request denied
						var out = {};
						out.cb_id = msg.cb_id;
						out.type = "reqd";
						out.error = "procedure with name '" + msg.procedure + "' is not defined";
						ws.send(JSON.stringify(out));

						// return, not going to run request
						return;
					}

					// send request accepted
					var out = {};
					out.cb_id = msg.cb_id;
					out.type = "reqa";
					ws.send(JSON.stringify(out));

					// call procedure
					var deferred = $q.defer();
					proc(msg.data, deferred);
					deferred.promise.then(function(data){
						sendPromiseFullfillment("res", msg.def_id, data);
					}, function(data) {
						sendPromiseFullfillment("rej", msg.def_id, data);
					}, function(data) {
						sendPromiseFullfillment("not", msg.def_id, data);
					});

					// all done
					break;

				case "reqa":
					// reqa is ignored here
					break;
			
				case "reqd":
					var deferred = AngoService.callbacks[msg.cb_id];
					delete AngoService.callbacks[msg.cb_id];
					deferred.reject(msg.err);
					break;

				case "res":
					var deferred = AngoService.callbacks[msg.cb_id];
					delete AngoService.callbacks[msg.cb_id];
					//++ TODO: see if lo_id is set instead of data, use linked object
					deferred.resolve(msg.data);
					break;

				case "rej":
					var deferred = AngoService.callbacks[msg.cb_id];
					delete AngoService.callbacks[msg.cb_id];
					//++ TODO: see if lo_id is set instead of data, use linked object
					deferred.reject(msg.data);
					break;

				case "not":
					var deferred = AngoService.callbacks[msg.cb_id];
					//++ TODO: see if lo_id is set instead of data, use linked object
					deferred.notify(msg.data);
					break;

				case "lor":
					AngoService.linkedObjectsCounter++;
					var id = AngoService.linkedObjectsCounter;
					AngoService.linkedObjects[id] = msg.data;
					var out = {
						type: "lora",
						lo_id: id,
						cb_id: msg.cb_id,
					};
					ws.send(JSON.stringify(out));
					break;

				case "lou":
					$apply(function() {
						// get linked object from register
						var obj = AngoService.linkedObjects[msg.lo_id];
						// remove old elements from registered linked object
						for (var key in obj) {
							delete obj[key];
						}
						// copy elements to registered linked object
						for (var key in msg.data) {
							obj[key] = msg.data[key];
						}
					});
					break;

				default:
					console.error("Unknown message type '" + msg.type +"'.");
					break;
			}
		}

		// all done
		return AngoService;
	});