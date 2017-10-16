// this file provides api for easy creation of charts

// chartData = [{x: unix time in miliseconds, y: number of lands/clicks}

function CreateChart(canvasElement, ChartData) {
    var chart = new Chart(canvasElement, {
        type: 'line',
        data: {
            datasets: [{
                label: '',
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                borderColor: 'rgb(54, 162, 235)',
                pointBorderWidth: 0,
                pointRadius: 0,
                hitRadius: 10,
                hoverRadius: 10,
                pointHoverBackgroundColor: 'rgb(54, 162, 235)',
                data: ChartData
            }]
        },
        options: {
            legend: {
                display: false
            },
            animation: false,
            scales: {
                xAxes: [{
                    type: 'time',
                    unit: '',
                    time: {
                        displayFormats: {
                            'millisecond': 'YYYY-MM-DD',
                            'second': 'YYYY-MM-DD',
                            'minute': 'YYYY-MM-DD',
                            'hour': 'YYYY-MM-DD',
                            'day': 'YYYY-MM-DD',
                            'week': 'YYYY-MM-DD',
                            'month': 'YYYY-MM-DD',
                            'quarter': 'YYYY-MM-DD',
                            'year': 'YYYY-MM-DD'
                        },
                        tooltipFormat: "YYYY-MM-DD",
                    },
                    position: 'bottom',
                }],
                yAxes: [{
                    display: true,
                    ticks: {
                        beginAtZero: true,
                    }
                }]
            }
        }
    });

    return chart
}