
// TODO: Handle service down
// TODO: Need to handle no dates (no scheduled events)
// TODO: Move to SSL for API call
// TODO: Move to DOM elements already existing in page (loadMap) and hookup callback from page

////////////////////////////////////////////
//
// Extend the Date object with some nice helpers
//
Date.prototype.displayTime = function() {
    var me = this;
    return function() {
        var hour = me.getHours();
        var dd = " AM";
        if (hour > 12) {
            hour = hour - 12;
            dd = " PM";
        }

        return " " + hour + ":" + ("0" + me.getMinutes()).slice(-2) + dd;
    }();
};

Date.prototype.displayDate = function() {
    var me = this;
    return function() {
        return (me.getMonth() + 1) + "/" + ("0" + me.getDate()).slice(-2) + "/" + me.getFullYear();
    }();
};

Date.prototype.justDate = function() {
    var me = this;
    return function() {
        return new Date(me.getFullYear(), me.getMonth(), me.getDate());
    }();
};


var Teams = (function() {

    ////////////////////////////////////////////
    //
    // Locations is a collection of Location objects
    //
    var Locations = (function() {
        var locations = {};

        var _key = function(lat, lng) {
            return lat + "|" + lng;
        };

        return {
            // Adds a location if it doesn't already exist
            addLocation: function(lat, lng, address, name) {
                var l = locations[_key(lat,lng)];
                if (l == undefined) {
                    // create the location
                    l = new Location(lat, lng, address, name);
                    locations[_key(lat,lng)] = l;
                }

                return l;
            },
            // Get all markers for the location for the specified date
            markers: function(dateObj) {
                var markers = [];
                for(var key in locations) {
                    var location = locations[key];
                    var marker = location.marker(dateObj);
                    if (marker != undefined) {

                        // Add the count to the marker icon
                        var len = marker.games.length;
                        if (len > 1) {
                            marker.setIcon('http://chart.apis.google.com/chart?chst=d_map_pin_letter&chld=' + len + '|FE7569');
                        }
                        markers.push(marker);
                    }
                }
                return markers;
            }
        }
    })();

    ////////////////////////////////////////////
    //
    // Location encapsulates attributes, games, and a map marker
    //
    var Location = (function(lat, lng, address, name) {
        var _lat = lat;
        var _lng = lng;
        var _address = address;
        var _name = name;
        var _markers = {};

        var _getInfo = function(marker) {

            var info = marker.info + "<ul>";

            marker.games.sort(function(game1, game2) {
                return game1.compare(game2);
            });

            marker.games.forEach(function (game) {
                info += "<li>" + game.toString() + "</li>"
            });
            info += "</ul>";

            return info;
        };

        var _key = function(dateObj) {
            return dateObj.displayDate();
        };

        return {
            // Adds a game to the current location
            addGame: function(dateObj, teamName) {
                // check if the marker exists
                var key = _key(dateObj);
                var marker = _markers[key];
                if (marker == undefined) {
                    // create the marker
                    marker = new google.maps.Marker({
                        position: {lat: _lat, lng: _lng},
                        title: _name || _address,
                        info: "<p>" + _name + "<br />" + _address + "</p>",
                        games: []
                    });
                    _markers[key] = marker;

                    // create the click handler
                    google.maps.event.addListener(marker, 'click', (function (marker) {
                        return function () {
                            _infoWindow.setContent(_getInfo(marker));
                            _infoWindow.open(map, marker);
                        }
                    })(marker));
                }

                marker.games.push(new Game(dateObj, teamName));
            },
            // Retrieves a marker for a specified date
            marker: function(dateObj) {
                return _markers[_key(dateObj)];
            }
        };
    });

    ////////////////////////////////////////////
    //
    // Game encapsulates game attributes and provides helpers to sort and convert to a display string
    //
    var Game = (function(dateObj, teamName) {
        var _dateObj = dateObj;
        var _teamName = teamName;

        return {
            compareString: function() { return _dateObj + _teamName; },
            toString: function() { return _dateObj.displayTime() + ": " + _teamName; },
            compare: function(game) {
                var c1 = this.compareString();
                var c2 = game.compareString();
                if (c1 > c2) return 1;
                if (c2 > c1) return -1;
                return 0;
            }
        };
    });

    var _infoWindow;
    var _datePicker; // DOM object
    var _map;  // Google Map object
    var _locations = Locations; // All locations

    function initMap() {

        _infoWindow = new google.maps.InfoWindow();
        _map = new google.maps.Map(document.getElementById('map'));
        _datePicker = document.getElementById('map-date-picker');

        // load team data
        getTeamsJSON(function (data) {
            processData(data);
        }, function (status) {
            // TODO: Handle this error differently
            alert('Faield to get data. Status: ' + status);
        });
    }

    var getTeamsJSON = function (successHandler, errorHandler) {

        var xhr = new XMLHttpRequest();
        var data;
        var status;

        xhr.open('Get', 'http://api.centralmarinsoccer.com/teams/');
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                status = xhr.status;
                if (status == 200) {
                    data = JSON.parse(xhr.responseText);
                    successHandler && successHandler(data);
                } else {
                    errorHandler && errorHandler(status);
                }
            }
        };

        xhr.send();
    };

    var processData = function (data) {

        var dates = {};

        // construct locations
        data.teams.forEach(function (team) {
            if (team.events) {
                team.events.forEach(function (event) {

                    var lat = event.location.latitude;
                    var lng = event.location.longitude;
                    if (!isNaN(lat) && !isNaN(lng)) {

                        var eventDate = new Date(event.start);
                        var location = _locations.addLocation(lat,lng, event.location.address, event.location.name);
                        location.addGame(eventDate, team.year + " " + team.gender + " " + team.level + " vs " + event.opponent);

                        // save off date
                        dates[eventDate.justDate()] = true;
                    }
                })
            }
        });

        var availableDates = [];
        for(var date in dates) {
            availableDates.push(new Date(date));
        }
        availableDates.sort(function(date1, date2) {
            // compare years first
            if (date1 > date2) return 1;
            if (date1 < date2) return -1;
            return 0;
        });
        availableDates.forEach(function(item) {
            var option = document.createElement('option');
            option.text = item.displayDate();
            option.value = item;
            _datePicker.add(option);
        });

        _datePicker.addEventListener("change", function() {
            updateMap(new Date(this.value));
        });

        updateMap(availableDates[0]);
    };

    var activeMarkers = [];
    var updateMap = function (dateObj) {

        // Clear existing markers
        activeMarkers.forEach(function(marker) {
            marker.setMap(null);
        });
        activeMarkers = [];

        // Add markers for the specified date
        var newBoundary = new google.maps.LatLngBounds();
        var markers = _locations.markers(dateObj);
        markers.forEach(function(marker) {
            marker.setMap(_map);
            newBoundary.extend(marker.position);

            activeMarkers.push(marker);
        });
        _map.fitBounds(newBoundary);
        var markerCluster = new MarkerClusterer(_map, markers, {imagePath: 'team/client/images/m'});
    };

    return {
        initMap : initMap
    }
})();

// function setupPage() {
//
//     var datePicker = document.createElement('select');
//     datePicker.id = "map-date-picker";
//     script.parentNode.appendChild(datePicker);
//
//     var div = document.createElement('div');
//     div.id = "map";
//     div.height = '400px';
//     div.width = '100%';
//     script.parentNode.appendChild(div);
// }
//
// setupPage();



