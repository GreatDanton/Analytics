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
                pointRadius: 3,
                pointBackgroundColor: 'rgb(54, 162, 235)',
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
            elements: {
                line: {
                    tension: 0 // disable smoothing => straight lines
                }
            },
            scales: {
                xAxes: [{
                    type: 'time',
                    ticks: {
                        autoSkip: true,
                    },
                    time: {
                        unit: 'day',
                        displayFormats: {
                            'millisecond': 'MM-DD',
                            'second': 'MM-DD',
                            'minute': 'MM-DD',
                            'hour': 'MM-DD',
                            'day': 'MM-DD',
                            'week': 'MM-DD',
                            'month': 'MM-DD',
                            'quarter': 'MM-DD',
                            'year': 'MM-DD'
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