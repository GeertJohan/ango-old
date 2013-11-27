
var app = angular.module("exampleApp", ["ango"]);

app.controller("exampleController", function($q, Ango) {
	Ango.RegisterProcedure("echo", function(data, deferred) {
		deferred.resolve(data);
	});
});