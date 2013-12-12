
var app = angular.module("exampleApp", ["ango"]);

app.controller("exampleController", function($q, Ango) {

	//++ controller can be destructed, what happens with this procedure then?
	Ango.RegisterProcedure("echo", function(data, deferred) {
		deferred.resolve(data);
	});

	Ango.Fire("")

	// make a call
	Ango.Call("getTime").then(function(time) {
		$scope.time = time;
	}, function(error) {
		console.error("reject: " + error);
	});
});