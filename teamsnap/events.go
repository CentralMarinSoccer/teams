package teamsnap

import (
	"time"
	"fmt"
)

func (ts TeamSnap) events(links relHrefDatas) []TeamEvent {

	var events []TeamEvent

	// Load all of the team event locations
	ts.team_locations(links)

	// Load the events
	if href, ok := links.findRelLink("events"); ok {
		tr, _ := ts.makeRequest(href)
		for _, e := range tr.Collection.Items {
			if event, ok := ts.event(e, ts.locations); ok {
				events = append(events, event)
			}
		}
	}

	return events
}

func (ts TeamSnap) event(e relHrefData, locs map[string]TeamEventLocation) (TeamEvent, bool) {

	if results, ok := e.Data.findValues("is_game", "name", "arrival_date", "duration_in_minutes", "division_location_id", "location_id", "minutes_to_arrive_early"); ok {
		if results["is_game"] != "true" {
			return TeamEvent{}, false
		}
		var loc TeamEventLocation
		locId := results["location_id"]
		if locId == "" {
			locId = results["division_location_id"]
		}
		loc = locs[locId]
		start, _ := time.Parse(time.RFC3339, results["arrival_date"])

		// Game start is arrival_date + minutes_to_arrive_early
		if (results["minutes_to_arrive_early"] != "") {
			if earlyArrival, err := time.ParseDuration(fmt.Sprintf("%sm", results["minutes_to_arrive_early"])); err == nil {
				start = start.Add(earlyArrival);
			}
		}

		// Only add events if they're for today or the future
		diff := start.Sub(time.Now())
		if diff.Hours() > -16 {

			var event =  TeamEvent{
				Start:    start,
				Opponent: ts.opponent(e.Links),
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

func (ts TeamSnap) team_locations(links relHrefDatas) {

	// Load club and division locations
	ts.load_locations(links, "division_locations")
	ts.load_locations(links, "locations")
}

func (ts TeamSnap) load_locations(links relHrefDatas, loc_type string) {

	if href, ok := links.findRelLink(loc_type); ok {
		tr, _ := ts.makeRequest(href)
		for _, l := range tr.Collection.Items {
			if results, ok := l.Data.findValues("id", "name", "address"); ok {
				ts.locations[results["id"]] = TeamEventLocation{
					Name:    results["name"],
					Address: results["address"],
				}
			}
		}
	}
}

func (ts TeamSnap) opponent(links relHrefDatas) string {
	if href, ok := links.findRelLink("opponent"); ok {
		tr, _ := ts.makeRequest(href)
		if results, ok := tr.Collection.Items[0].Data.findValues("name"); ok {
			return results["name"]
		}
	}

	return ""
}

