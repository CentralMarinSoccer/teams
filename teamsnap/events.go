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

	log.Debugf("Values: %v", results)

	// Use the appropriate location
	var location TeamEventLocation
	if results["location_id"] != "" {
		log.Debugf("Location ID: %s - Locations: %v", results["location_id"], ts.locations)

		location = ts.locations[results["location_id"]]
	} else {
		log.Debugf("Division Location ID: %s - Locations: %v", results["division_location_id"], ts.divisionLocations)

		location = ts.divisionLocations[results["division_location_id"]]
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

		log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Opponent id: %s, opponent map: %v", results["opponent_id"], ts.opponents)

		var event = TeamEvent{
			Start:    start,
			Opponent: ts.opponents[results["opponent_id"]],
			Duration: results["duration_in_minutes"],
			Location: TeamEventLocation{
				Name:    location.Name,
				Address: location.Address,
			},
		}

		log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Event before geocoding: %v", event)

		if location.Address != "" {
			// Geocode address to latitude / longitude
			if address := ts.configuration.Geocoder.Lookup(location.Address); address != nil {
				event.Location.Address = address.FormattedAddress
				event.Location.Latitude = address.Lat
				event.Location.Longitude = address.Lng

				log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Successfully geocoded address: %s for %s", location.Address, location.Name)
			} else {
				log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unable to geocode address: %s for %s", location.Address, location.Name)
			}
		} else {
			log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Skipping geocoding address: %s for %s", location.Address, location.Name)
		}

		log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Event after geocoding: %v", event)

		team.Events = append(team.Events, event)
	}

	return true
}
