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

        var chartdata1 = {
            chart: { type: 'line' },
            title: { text: 'Number of nodes' },
            xAxis: { categories: ['20', '30', '40', '50'] },
            yAxis: { title: { text: 'Time (ms)' } },
            credits: { enabled: false },
            plotOptions: { line: { dataLabels: { enabled: true },  enableMouseTracking: false } },
            series: []
        };
        
        if ( p == "block_valid_ratio_percentage") {
            chartdata1.yAxis.title.text = 'Percentage (%)';
        }

        // var timeouts = [200, 300];
        // var speeds = [2, 5];

        console.debug("Working on 200 and 2");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/200/2" )
            .then(function (response) {
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_200_2', chartdata1);
            });

        console.debug("Working on 200 and 5");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/200/5" )
            .then(function (response) {
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_200_5', chartdata1);
            });

        console.debug("Working on 300 and 2");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/300/2" )
            .then(function (response) {
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_300_2', chartdata1);
            });

        console.debug("Working on 300 and 5");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/300/5" )
            .then(function (response) {
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_300_5', chartdata1);
            });

        $('#refreshButton').html('Refresh');
        $('#refreshButton').prop("disabled", false);

    };

    //initial load
    $scope.loadData();
});