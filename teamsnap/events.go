package teamsnap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"time"
)

const pastGameTime = -2

func processEvent(results nameValueResults, ts TeamSnap, team *Team) bool {

	// We only care about games
	if results["is_game"] != "true" {
		return true
	}

	// Use the appropriate location
	var location TeamEventLocation
	if results["location_id"] == "" {
		location = ts.locations[results["location_id"]]
	} else {
		location = ts.divisionLocations[results["location_id"]]
	}

	// Game start is arrival_date + minutes_to_arrive_early
	start, _ := time.Parse(time.RFC3339, results["arrival_date"])
	if results["minutes_to_arrive_early"] != "" {
		if earlyArrival, err := time.ParseDuration(fmt.Sprintf("%sm", results["minutes_to_arrive_early"])); err == nil {
			start = start.Add(earlyArrival)
		}
	}

	// Only add events if they're for today or the future
	diff := start.Sub(time.Now())
	if diff.Hours() > pastGameTime {

		var event = TeamEvent{
			Start:    start,
			Opponent: ts.opponents[results["opponent_id"]],
			Duration: results["duration_in_minutes"],
			Location: TeamEventLocation{
				Name:    location.Name,
				Address: location.Address,
			},
		}

		if location.Address != "" {
			// Geocode address to latitude / longitude
			if address := ts.configuration.Geocoder.Lookup(location.Address); address != nil {
				event.Location.Address = address.FormattedAddress
				event.Location.Latitude = address.Lat
				event.Location.Longitude = address.Lng
			} else {
				log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unable to geocode address: %s for %s", location.Address, location.Name)
			}
		}

		team.Events = append(team.Events, event)
	}

	return true
}
