var elements = ["main", "pollutionValue",
    "pollutionType",
    "pollutionPopup", "stationName", "chartCont",
    "airQuality", "selectPollution", "chart", "pollutionCBox"
];
elements.forEach(el => {
    var el = document.getElementById(el);
});


var map
var geocoder

function init() {
    geocoder = new google.maps.Geocoder();
    map = new google.maps.Map(document.getElementById('map'), {
        center: {
            lat: 52.2330653,
            lng: 20.9211125
        },
        zoom: 6
    });
}
var markers = [];

var searchInput = document.getElementById("mapInput");
searchInput.addEventListener("keyup", function (event) {
    event.preventDefault();
    if (event.keyCode === 13) {
        Search();
    }
});

function Search() {
    var address = document.getElementById('mapInput').value;
    geocoder.geocode({
        'address': address
    }, function (results, status) {
        if (status == 'OK') {
            map.setCenter(results[0].geometry.location);
            map.setZoom(11)
        } else {
            alert('Nie znaleziono miejsca');
        }
    });
}

function toggleStationInfo() {
    if (pollutionPopup.classList.contains("block--none")) {
        pollutionPopup.classList.remove("block--none");
        main.classList.add("block--blur");
    } else {
        pollutionPopup.classList.add("block--none");
        main.classList.remove("block--blur");
    }
}

var pChart = null;

function postDataToGraph(id) {
    fetch('http://api.gios.gov.pl/pjp-api/rest/data/getData/' + id, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            pollutionType.innerHTML = jsonData.key;
            if (jsonData.values[0].value != null) {
                pollutionValue.innerHTML = Math.floor(jsonData.values[0].value) + ' µg/m3'
            } else {
                pollutionValue.innerHTML = 'brak danych'
            }

            if (pChart != null) {
                pChart.destroy();
            }

            var labels = jsonData.values.map(function (e) {
                return e.date.split(' ')[1].split(':')[0];
            });
            var data = jsonData.values.map(function (e) {
                return Math.floor(e.value);
            });


            var ctx = chart.getContext("2d");
            pChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Poziom zanieczyszenia o danej godzinie',
                        data: data
                    }]
                }
            })


        })
        .catch((err) => {
            console.log(err.message);
        });
}

function showStationInfo(id, name) {
    fetch('http://api.gios.gov.pl/pjp-api/rest/aqindex/getIndex/' + id, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            stationName.innerHTML = name;
            airQuality.classList.add('color--' + jsonData.stIndexLevel.id)
            airQuality.innerHTML = jsonData.stIndexLevel.indexLevelName;
            toggleStationInfo();
        })
        .catch((err) => {
            console.log(err.message);
        });

    fetch('http://api.gios.gov.pl/pjp-api/rest/station/sensors/' + id, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            postDataToGraph(jsonData[0].id)
            jsonData.forEach((el) => {
                var option = document.createElement("option");
                option.value = el.id;
                option.text = el.param.paramName;
                selectPollution.add(option);
            })
        })
        .catch((err) => {
            console.log(err.message);
        });
}

function getStationStatus(sID) {
    var pID
    fetch('http://api.gios.gov.pl/pjp-api/rest/aqindex/getIndex/' + sID, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            pID = jsonData.stIndexLevel.id;
        })
        .catch((err) => {
            console.log(err.message);
        });
    return pID;
}

function getImage(type, name) {
    var link = 'static/img/';
    switch (type) {
        case 'pollution':
            switch (getStationStatus(name.id)) {
                case '0':
                    link += 'pollution/0.png';
                    break;
                case '1':
                    link += 'pollution/1.png';
                    break;
                case 2:
                    link += 'pollution/2.png';
                    break;
                case 3:
                    link += 'pollution/3.png';
                    break;
                case 4:
                    link += 'pollution/4.png';
                    break;
                case 5:
                    link += 'pollution/5.png';
                    break;
                case 'nodata':
                    link += 'pollution/00.png';
                    break;
                default:
                    link += 'pollution/00.png';
            }
            break;
        default:
    }
    console.log(link)
    return link
}

function showPollutionMarkers() {
    pollutionCBox.classList.remove('show-el-checkbox--no');
    pollutionCBox.classList.add('show-el-checkbox--yes');
    pollutionCBox.style = 'pointer-events:none;';

    const pollutionUrl = 'http://api.gios.gov.pl/pjp-api/rest/station/findAll';

    fetch(pollutionUrl, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            jsonData.forEach((obj) => {
                var image = {
                    url: getImage('pollution', obj),
                    scaledSize: new google.maps.Size(32, 32),
                    origin: new google.maps.Point(0, 0),
                    anchor: new google.maps.Point(0, 0)
                };
                var marker = new google.maps.Marker({
                    position: new google.maps.LatLng(obj.gegrLat, obj.gegrLon),
                    map: map,
                    title: obj.stationName,
                    icon: image
                });
                google.maps.event.addListener(marker, 'click', function () {
                    showStationInfo(obj.id, obj.stationName)
                });
                markers.push(marker);
            });
            var markerCluster = new MarkerClusterer(map, markers, {
                imagePath: 'https://developers.google.com/maps/documentation/javascript/examples/markerclusterer/m'
            });
        })
        .catch((err) => {
            console.log(err.message);
        });
    var markerCluster = new MarkerClusterer(map, markers);
}

function showPharmacyInfo(id) {
    alert(id)
}

function showPharmaciesMarkers() {
    pharmaciesCBox.classList.remove('show-el-checkbox--no');
    pharmaciesCBox.classList.add('show-el-checkbox--yes');
    pharmaciesCBox.style = 'pointer-events:none;';

    const pharmaciesUrl = 'static/json/pharmacies-markers.json';

    fetch(pharmaciesUrl, {
            method: 'GET',
            mode: 'cors'
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
        })
        .then((jsonData) => {
            jsonData.PharmacyMarkers.forEach((obj) => {
                var image = {
                    url: obj.Image,
                    scaledSize: new google.maps.Size(32, 32),
                    origin: new google.maps.Point(0, 0),
                    anchor: new google.maps.Point(0, 0)
                };
                var marker = new google.maps.Marker({
                    position: new google.maps.LatLng(obj.Lat, obj.Lng),
                    map: map,
                    icon: image
                });
                google.maps.event.addListener(marker, 'click', function () {
                    showPharmacyInfo(obj.ID)
                });
                var w
                for (i = 0; i < markers.length; i++) {
                    if (markers[i].position == marker.position) {
                        w++;
                    }
                }
                if (w == 0) {
                    markers.push(marker);
                }
            });
            var markerCluster = new MarkerClusterer(map, markers, {
                imagePath: 'https://developers.google.com/maps/documentation/javascript/examples/markerclusterer/m'
            });
        })
        .catch((err) => {
            console.log(err.message);
        });
    var markerCluster = new MarkerClusterer(map, markers);
}

var info = [{
    name: "pharmacies",
    symbol: "",
    fullname: "Apteki",
    content: ""
}, {
    name: "pollution",
    symbol: "",
    fullname: "Zanieczyszenie",
    content: "Ta funkcja pozwala Ci przeglądać dane dotyczące zanieczyszenia w Polsce. Aby ją włączyć kliknij przycisk po prawej stronie pytajnika. Zobaczysz mapę z różnokolorowymi punktami, które oznaczają stopnie zanieczyszenia (od zielonego do czerownego). Kliknięcie na znaczniki da ci natomiast możliwość zobaczenia dokładnych wartości obecnych i historycznych danych zanieczyszczeń w powietrzu."
}]

function hideInfo() {
    main.classList.remove("block--blur");
    infoPopup.classList.add("block--none");
}

function showInfo(name) {
    info.forEach((el) => {
        if (el.name == name) {
            infoTitle.innerHTML = el.fullname;
            infoText.innerHTML = el.content;
            infoIcon.innerHTML = el.symbol;
            main.classList.add("block--blur");
            infoPopup.classList.remove("block--none");
        }
    })
}