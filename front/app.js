console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var cp = $('#selectProperty').val();

        var cp_array = str.split(" - ");
        var c = str.trim(cp_array[0]);
        var p = str.trim(cp_array[1]);

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