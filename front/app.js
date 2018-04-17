console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var cp = $('#selectProperty').val();

        var cp_array = cp.split(" - ");
        var c = cp_array[0].trim();
        var p = cp_array[1].trim();

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