
angular.module('ango', [])
	.factory('Ango', function($q, $http) {
		var service = {};

		// Send is a fire-and-forget method to send data to a service
		// This looks up a normal (no callback) service handler
		service.Fire = function(service, data) {
			//++
			return;
		}

		// Call sends a request and returns a promise.
		// The returned deferred is resolved by the first response from server.
		// Also allows server to send notifications (e.g. processing progress).
		// The given header
		service.Call = function(service, data) {
			var deferred = $q.defer();
			//++
			return deferred.promise;
		};

		// RegisterProcedure allows you to implement client-side procedure that can be called by the server.
		// When the handler returns a promise and the server expects a callback the resolve/reject is sent to server
		// When the server expects a callback, but no promise is returns, an error is sent to the server.
		service.RegisterProcedure = function(name, handler) {
			if(typeof name != "string") {
				throw("Invalid procedure name, not a string");
			}
			if(typeof handler != "function") {
				throw("Invalid procedure handler, not a function");
			}

			service.procedures[name] = handler;
		}

		//++ open websocket
		//++ listen for messages, parse them and switch by type
		//++ on "req", create deferred, call handler with (data, deferred)
		//++ check if deferred was done, otherwise reject with "procedure handler did not resolve or reject"

		// all done
		return service;
	});