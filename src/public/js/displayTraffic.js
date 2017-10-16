// displaying user lands traffic
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

// displaying clicks
var clicksTraffic = JSON.parse(document.getElementById('user-clicks').innerText);
var clicksCard = document.getElementById('numOfClicks');

// display number of clicks in small info card
var numOfClicks = clicksTraffic['NumOfClicks'];
if (numOfClicks == "") {
    clicksCard.innerHTML = 0;
} else {
    clicksCard.innerHTML = numOfClicks;
}

var landChart = document.getElementById('landsChart');
var clicksChart = document.getElementById('clicksChart');

// creates land chart data from parsed json coming
// from server db
function createLandChartData(landTraffic) {
    var arr = [];
    var json = landTraffic['Lands'];
    // if json data does not exist return empty array
    if (json == null) {
        return arr;
    }

    for (var element of landTraffic['Lands']) {
        var obj = {}
        obj['x'] = element.Date
        obj['y'] = element.LandNumber;
        arr.push(obj)
    }
    return arr;
}

// creates clicks chart data parsed from json
// coming from server
function createClicksChartData(clicksTraffic) {
    var arr = [];
    var json = clicksTraffic['Clicks'];
    // if json data does not exist return empty arr
    if (json == null) {
        return arr
    }
    for (var element of json) {
        var obj = {}
        obj['x'] = element.Date
        obj['y'] = element.ClicksNum;
        arr.push(obj)
    }
    return arr
}

var landData = createLandChartData(landTraffic);
var clicksData = createClicksChartData(clicksTraffic);
CreateChart(landChart, landData);
CreateChart(clicksChart, clicksData);