package teamsnap

import (
	"time"
	"fmt"
)

//type location struct {
//	name    string
//	address string
//}

func (ts TeamSnap) events(links relHrefDatas) []TeamEvent {

	var events []TeamEvent

	// Load all of the team event locations
	loc := ts.locations(links)

	// Load the events
	if href, ok := links.findRelLink("events"); ok {
		tr, _ := ts.makeRequest(href)
		for _, e := range tr.Collection.Items {
			if event, ok := ts.event(e, loc); ok {
				events = append(events, event)
			}
		}
	}

	return events
}

func (ts TeamSnap) event(e relHrefData, locs map[string]TeamEventLocation) (TeamEvent, bool) {
	if results, ok := e.Data.findValues("type", "name", "arrival_date", "duration_in_minutes", "location_id", "minutes_to_arrive_early"); ok {
		loc := locs[results["location_id"]]
		start, _ := time.Parse(time.RFC3339, results["arrival_date"])

		// Game start is arrival_date + minutes_to_arrive_early
		if earlyArrival, err := time.ParseDuration(fmt.Sprintf("%sm", results["minutes_to_arrive_early"])); err == nil {
			start = start.Add(earlyArrival);
		}

		// Only add events if they're for today or the future
		diff := start.Sub(time.Now())
		if diff.Hours() > -16 {

			var event =  TeamEvent{
				Name:     results["name"],
				Start:    start,
				Duration: results["duration_in_minutes"],
				Location: TeamEventLocation{
					Name:    loc.Name,
					Address: loc.Address,
				},
			}

			// Geocode address to latitude / longitude
			if address := ts.configuration.Geocoder.Lookup(loc.Address); address != nil {
				event.Location.Address = address.FormattedAddress
				event.Location.Latitude = address.Lat
				event.Location.Longitude = address.Lng
			}

			return event, true
		}
	}

	return TeamEvent{}, false
}

func (ts TeamSnap) locations(links relHrefDatas) map[string]TeamEventLocation {
	locs := make(map[string]TeamEventLocation)

	// Load all of the locations for this team
	if href, ok := links.findRelLink("locations"); ok {
		tr, _ := ts.makeRequest(href)
		for _, l := range tr.Collection.Items {
			if results, ok := l.Data.findValues("id", "name", "address"); ok {
				locs[results["id"]] = TeamEventLocation{
					Name:    results["name"],
					Address: results["address"],
				}

			}
		}
	}

	return locs
}