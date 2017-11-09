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

// changes Count, Date from JSON object to {x, y}
// ie. format that chartJS understand
function prepareData(trafficArr) {
    arr = [];
    for (i = 0; i < trafficArr.length; i++) {
        obj = {};
        obj['x'] = trafficArr[i].Date;
        obj['y'] = trafficArr[i].Count;
        arr.push(obj);
    }
    return arr;
}

var landTraffic = displayLandTraffic()
var clicksTraffic = displayClicks()

var landChart = document.getElementById('landsChart');
var clicksChart = document.getElementById('clicksChart');

var landData = prepareData(landTraffic["Lands"]);
var clickData = prepareData(clicksTraffic["Clicks"]);

CreateChart(landChart, landData);
CreateChart(clicksChart, clickData);