
// TODO: Handle service down
// TODO: Need to handle no dates (no scheduled events)

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
                            marker.setIcon('https://chart.googleapis.com/chart?chst=d_map_pin_letter&chld=' + len + '|FE7569');
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

    var _infoWindow;  // Google Info popup window
    var _datePicker; // DOM object
    var _map;  // Google Map object
    var _locations = Locations; // All locations
    var _teams  = []

    var teamsFn = doT.template("<table><tr><td><ul>{{~it :value:index}} <li data-index='{{=index}}'>{{=value.year}} {{=value.gender}} {{=value.level}} - {{=value.name}}</li> {{~}}</ul></td><td valign='top'><div id='team'/></td></tr></table>");
    var teamFn = doT.template("<ul>{{~it :value:index}} <li data-index='{{=index}}'>{{=value.name}} {{=value.type}}</li> {{~}}");

    var teamName = function(team) {
        return team.year + " " + team.gender + " " + team.level;
    };

    // Create the XHR object.
    function createCORSRequest(method, url) {
      var xhr = new XMLHttpRequest();
      if ("withCredentials" in xhr) {
        // XHR for Chrome/Firefox/Opera/Safari.
        xhr.open(method, url, true);
      } else if (typeof XDomainRequest != "undefined") {
        // XDomainRequest for IE.
        xhr = new XDomainRequest();
        xhr.open(method, url);
      } else {
        // CORS not supported.
        xhr = null;
      }
      return xhr;
    }

    var initMap = function() {

        _infoWindow = new google.maps.InfoWindow();
        _map = new google.maps.Map(document.getElementById('map'));
        _datePicker = document.getElementById('map-date-picker');

        // load team data
        getTeamsJSON(function (data) {
            processData(data);
        }, function (status) {
            // TODO: Handle this error differently
            alert('Failed to get data. Status: ' + status);
        });
    };

    var getTeamsJSON = function (successHandler, errorHandler) {

        var data;
        var status;
        var xhr = createCORSRequest('GET', 'https://api.centralmarinsoccer.com/teams/');
        if (!xhr) {
            alert('CORS not supported');
            return;
        }

        xhr.onload = function() {
            data = JSON.parse(xhr.responseText);
            successHandler && successHandler(data); 
        };

        xhr.onerror = function() {
          errorHandler && errorHandler(xhr.status);
        };

        xhr.send();
    };
    var processData = function (data) {

        var dates = {};

        // construct locations
        data.teams.forEach(function (team) {

            _teams.push(team);

            if (team.events) {
                team.events.forEach(function (event) {

                    var lat = event.location.latitude;
                    var lng = event.location.longitude;
                    if (!isNaN(lat) && !isNaN(lng)) {

                        var eventDate = new Date(event.start);
                        var location = _locations.addLocation(lat,lng, event.location.address, event.location.name);
                        location.addGame(eventDate, teamName(team) + " vs " + event.opponent);

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

        // Sort and Display teams
        _teams.sort(function(team1, team2) {
            if (team1.gender > team2.gender) return 1;
            if (team1.gender < team2.gender) return -1;

            if (team1.year > team2.year) return 1;
            if (team1.year < team2.year) return -1;

            if (team1.level > team2.level) return 1;
            if (team1.level < team2.level) return -1;

            return 0;
        });
        displayTeams();
    };

    var displayTeams = function() {
        var teams = document.getElementById("teams");
        teams.innerHTML = teamsFn(_teams);
	var _teamId = document.getElementById("team");

        teams.addEventListener('click', function (event) {
          var index = event.target.getAttribute('data-index');
          var team = _teams[index];

          team.members.sort(function(member1, member2) {
              // sort on type first
              if (member1.type > member2.type) return 1;
              if (member1.type < member2.type) return -1;

              if (member1.name > member2.name) return 1;
              if (member1.name < member2.name) return -1;

              return 0;
          });

          _teamId.innerHTML = teamFn(team.members);
        });
    };

    var activeMarkers = [];
    var markerCluster;
    var updateMap = function (dateObj) {

        if (markerCluster == undefined) {
            markerCluster = new MarkerClusterer(_map, null, {imagePath: 'https://api.centralmarinsoccer.com/teams/static/images/m'});
        }
        // Clear existing markers
        activeMarkers.forEach(function(marker) {
            marker.setMap(null);
        });
        activeMarkers = [];
        markerCluster.clearMarkers();

        // Add markers for the specified date
        var newBoundary = new google.maps.LatLngBounds();
        var markers = _locations.markers(dateObj);
        markers.forEach(function(marker) {
            marker.setMap(_map);
            newBoundary.extend(marker.position);

            activeMarkers.push(marker);
        });
        _map.fitBounds(newBoundary);
        markerCluster.addMarkers(markers);
    };

    return {
        initMap : initMap
    }
})();



