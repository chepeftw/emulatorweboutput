console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {
    $http.get("http://api.chepeftw.com/property/sent_messages/messages_count")
        .then(function (response) {

            $scope.sortType = 'name'; // set the default sort type
            $scope.sortReverse = false;  // set the default sort order
            $scope.searchFish = '';     // set the default search/filter term

            $scope.sushi = response.data;
        });
});