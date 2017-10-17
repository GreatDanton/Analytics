// displaying user lands traffic
function displayLandTraffic() {
    var landTraffic = JSON.parse(document.getElementById('user-lands-traffic').innerText);
    var landsCard = document.getElementById('numOfLands');

    // display number of lands in small info card
    var numOfLands = landTraffic['NumOfLands'];
    // if numOfLands is empty => no lands happened
    if (numOfLands == "") {
        landsCard.innerHTML = 0;
    } else {
        // display actual number of lands
        landsCard.innerHTML = numOfLands;
    }
    return landTraffic;
}

// displaying clicks
function displayClicks() {
    var clicksTraffic = JSON.parse(document.getElementById('user-clicks').innerText);
    var clicksCard = document.getElementById('numOfClicks');

    // display number of clicks in small info card
    var numOfClicks = clicksTraffic['NumOfClicks'];
    if (numOfClicks == "") {
        clicksCard.innerHTML = 0;
    } else {
        clicksCard.innerHTML = numOfClicks;
    }
    return clicksTraffic;
}

var landTraffic = displayLandTraffic()
var clicksTraffic = displayClicks()

var landChart = document.getElementById('landsChart');
var clicksChart = document.getElementById('clicksChart');


// preparesChartDate prepares chart data out of json string
//
function prepareChartData(jsonData, nameOfClicks) {
    var arr = [];
    var data = jsonData;
    // if json data does not exist return empty array
    if (data == null) {
        return arr;
    }
    var dayMS = 24 * 60 * 60 * 1000; // miliseconds

    for (var i = 0; i < data.length; i++) {
        var date1 = data[i].Date;

        // last item in array
        if (i + 1 == data.length) {
            var obj = { 'x': date1, 'y': data[i][nameOfClicks] }
            arr.push(obj)
            break;
        }
        var date2 = data[i + 1].Date;

        var obj = {}
        obj['x'] = date1;
        obj['y'] = data[i][nameOfClicks];
        arr.push(obj);

        // if the date has no data, this makes sure
        // the number for that data is 0
        //(chartJS is not able to do that automatically)
        if (date1 + dayMS < date2) {
            d = date1 + dayMS
            while (d < date2) {
                var obj = { 'x': d, 'y': 0 }
                arr.push(obj)
                d += dayMS;
            }
        }
    }
    return arr;
}

//var landData = createLandChartData(landTraffic);
var landData = prepareChartData(landTraffic['Lands'], 'LandNumber');
var clicksData = prepareChartData(clicksTraffic['Clicks'], 'ClicksNum')
//var clicksData = createClicksChartData(clicksTraffic);
CreateChart(landChart, landData);
CreateChart(clicksChart, clicksData);