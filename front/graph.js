console.debug("Hi");
var app = angular.module('myApp', []);
app.controller('myCtrl', function ($scope, $http) {

    $scope.loadData = function () {
        var cp = $('#selectProperty').val();
        var fltr = $('#inputFilter').val();

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
            // plotOptions: { column: { stacking: 'normal', dataLabels: { enabled: true, color: (Highcharts.theme && Highcharts.theme.dataLabelsColor) || 'white' } }, line: { dataLabels: { enabled: true },  enableMouseTracking: false } },
            exporting: { enabled: true, filename: 'MyChart' },
            series: []
        };

        var graphFileName = ""
        if ( p == "query_complete_ms") {
            $scope.graphTitle = "Monitor Query Complete";
            graphFileName = 'QueryComplete';
            chartdata1.chart.type = 'line';
        } else if ( p == "block_valid_ratio_percentage") {
            $scope.graphTitle = "Block Valid Ratio";
            chartdata1.yAxis.title.text = 'Percentage (%)';
            graphFileName = 'Accuracy';
            chartdata1.chart.type = 'column';
            chartdata1.yAxis.min = 97;
        } else if ( c == "router_sent_messages" && p == "messages_count") {
            $scope.graphTitle = "Router Message Count";
            chartdata1.yAxis.title.text = 'Number of Packets';
            graphFileName = 'RouterMessageCount';
            chartdata1.chart.type = 'column';
            chartdata1.plotOptions = { column: { stacking: 'normal', dataLabels: { enabled: true, color: (Highcharts.theme && Highcharts.theme.dataLabelsColor) || 'white' } }, line: { dataLabels: { enabled: true },  enableMouseTracking: false } };
        } else if ( c == "raft_sent_messages" && p == "messages_count") {
            $scope.graphTitle = "Raft Message Count";
            chartdata1.yAxis.title.text = 'Number of Packets';
            graphFileName = 'RaftMessageCount';
            chartdata1.chart.type = 'column';
            chartdata1.plotOptions = { column: { stacking: 'normal', dataLabels: { enabled: true, color: (Highcharts.theme && Highcharts.theme.dataLabelsColor) || 'white' } }, line: { dataLabels: { enabled: true },  enableMouseTracking: false } };
        }

        var postfix = ""
        if (fltr) {
            postfix = "?filter=" + fltr;
        }

        console.debug("Working on 200 and 2");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/200/2" + postfix )
            .then(function (response) {
                chartdata1.exporting.filename = graphFileName + '200ms2ms';
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_200_2', chartdata1);
            });

        console.debug("Working on 200 and 5");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/200/5" + postfix )
            .then(function (response) {
                chartdata1.exporting.filename = graphFileName + '200ms5ms';
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_200_5', chartdata1);
            });

        console.debug("Working on 300 and 2");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/300/2" + postfix )
            .then(function (response) {
                chartdata1.exporting.filename = graphFileName + '300ms2ms';
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_300_2', chartdata1);
            });

        console.debug("Working on 300 and 5");
        $http.get( "http://api.chepeftw.com/graph/" + c + "/" + p + "/300/5" + postfix )
            .then(function (response) {
                chartdata1.exporting.filename = graphFileName + '300ms5ms';
                chartdata1.series = response.data.Highchart;
                Highcharts.chart('container_300_5', chartdata1);
            });

        $('#refreshButton').html('Refresh');
        $('#refreshButton').prop("disabled", false);

    };

    //initial load
    $scope.loadData();
});