console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var apiURL = "http://api.chepeftw.com/graph";

        var chartdata = {
            chart: {
                type: 'line'
            },
            title: {
                text: 'Number of nodes'
            },
            xAxis: {
                categories: ['20', '30', '40', '50']
            },
            yAxis: {
                title: {
                    text: 'Time (ms)'
                }
            },
            plotOptions: {
                line: {
                    dataLabels: {
                        enabled: true
                    },
                    enableMouseTracking: false
                }
            },
            series: []
        };

        $http.get( apiURL )
            .then(function (response) {
                // $scope.sushi = response.data;

                var seriesGD = [],
                    graphData= response.data.Highchart;

                for (var i=0; i< graphData.length; i++) {
                    seriesGD.push({"name" : graphData[i].name, "data" : graphData[i].data})
                }

                chartdata.series = seriesGD;

                Highcharts.chart('container', chartdata);

            });
    };

    //initial load
    $scope.loadData();
});