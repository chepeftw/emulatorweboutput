console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var apiURL = "http://api.chepeftw.com/graph/monitor_query_complete/query_complete_ms/200/2";

        var chartdata1 = {
            chart: { type: 'line' },
            title: { text: 'Number of nodes' },
            xAxis: { categories: ['20', '30', '40', '50'] },
            yAxis: { title: { text: 'Time (ms)' } },
            plotOptions: { line: { dataLabels: { enabled: true },  enableMouseTracking: false } },
            series: []
        };

        var timeouts = [200, 300];
        var speeds= [2, 5];

        var i, j;
        for (i = 0; i < timeouts.length; i++) {
            for (j = 0; i < speeds.length; j++) {
                var apiURL1 = "http://api.chepeftw.com/graph/monitor_query_complete/query_complete_ms/" + timeouts[i] + "/" + speeds[j];
                $http.get( apiURL1 )
                    .then(function (response) {
                        chartdata1.series = response.data.Highchart;
                        Highcharts.chart('container_mqc_' + timeouts[i] + '_' + speeds[j], chartdata1);
                    });
            }
        }

        var i, j;
        for (i = 0; i < timeouts.length; i++) {
            for (j = 0; i < speeds.length; j++) {
                var apiURL1 = "http://api.chepeftw.com/graph/monitor_accuracy/block_valid_ratio_percentage/" + timeouts[i] + "/" + speeds[j];
                $http.get( apiURL1 )
                    .then(function (response) {
                        chartdata1.series = response.data.Highchart;
                        Highcharts.chart('container_bvrp_' + timeouts[i] + '_' + speeds[j], chartdata1);
                    });
            }
        }

    };

    //initial load
    $scope.loadData();
});