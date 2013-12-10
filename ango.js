
angular.module('ango', [])
	.factory('Ango', function($q, $http) {
		var AngoService = {
			procedures: {},
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
			var deferred = $q.defer();
			//++
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

				case "res":
					//++ resolve outstanding deferred
					break;

				case "rej":
					//++ reject outstanding deferred
					break;

				case "not":
					//++ notify outstanding deferred
					break;

				case "lor":
					//++ register linked object
					break;

				case "lou":
					//++ update linked object
					break;

				default:
					console.error("Unknown message type '" + msg.type +"'.");
					break;
			}
		}

		// all done
		return AngoService;
	});