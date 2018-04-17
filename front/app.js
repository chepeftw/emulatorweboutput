console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var c = $('#selectCollection').val();
        var p = $('#selectProperty').val();

        $('#refreshButton').html('Loading...');
        $('#refreshButton').prop("disabled", true);

        // SelectElement("selectCollection", c);
        // SelectElement("selectProperty", p);

        $http.get("http://api.chepeftw.com/property/" + c + "/" + p)
            .then(function (response) {
                $scope.sushi = response.data;
                $('#refreshButton').html('Refresh');
                $('#refreshButton').prop("disabled", false);
            });
    };

    //initial load
    $scope.loadData();
});