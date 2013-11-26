
angular.module('ango', [])
	.factory('Ango', function($q, $http) {
		var service = {};

		// Send is a fire-and-forget method to send data to a service
		// This looks up a normal (no callback) service handler
		service.Send = function(service, data) {
			//++
			return;
		}

		// Request sends a request and returns a promise.
		// The returned deferred is resolved by the first response from server.
		// Also allows server to send notifications (e.g. processing progress).
		// The given header
		service.Request = function(service, data) {
			var deferred = $q.defer();
			//++
			return deferred.promise;
		};

		// RegisterService allows you to implement client-side services that can be called by the server.
		// When the handler returns a promise and the server expects a callback the resolve/reject is sent to server
		// When the server expects a callback, but no promise is returns, an error is sent to the server.
		service.RegisterService = function(name, handler) {
			//++
		}

		// all done
		return service;
	});