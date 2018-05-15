console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var apiURL = "http://api.chepeftw.com/graph";

        $http.get( apiURL )
            .then(function (response) {
                // $scope.sushi = response.data;

                var seriesGD = [],
                    graphData= response.data;

                for (var i=0; i< graphData.length; i++) {
                    seriesGD.push({"name" : graphData[i].name, "data" : graphData[i].data})
                }


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

                chartdata.series = seriesGD;

                Highcharts.chart('container', chartdata);

            });
    };

    //initial load
    $scope.loadData();
});