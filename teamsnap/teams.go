package teamsnap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/url"
	"strings"
)

// TODO: Move Division Location outside loop since it doesn't change from Team to Team
// TODO: Does order of bulk load types matter? We're using a map today, so we can't control order. We can add a sequence number and sort based on that if necessary.

type itemProcesserFunc func(results nameValueResults, ts TeamSnap, team *Team) bool

type itemType struct {
	itemNames     []string
	itemProcessor itemProcesserFunc
}

// Map bulk categories to values we care about
var categories = map[string]itemType{
	"member": {
		[]string{"first_name", "last_name", "birthday", "is_manager", "is_non_player", "is_owner", "position"},
		processMember,
	},
	"event": {
		[]string{"is_game", "name", "arrival_date", "duration_in_minutes", "division_location_id", "location_id", "minutes_to_arrive_early", "opponent_id"},
		processEvent,
	},
	"team_preferences": {
		[]string{"gender"},
		processTeamPreferences,
	},
	"team_photo": {
		nil,
		nil,
	},
	"opponent": {
		[]string{"id", "name"},
		nil,
	},
	"location": {
		[]string{"id", "name", "address"},
		nil,
	},
	"division_location": {
		[]string{"id", "name", "address"},
		nil,
	},
}

func (ts TeamSnap) teams() (Teams, bool) {

	// If root not loaded, load it
	if ts.root.Collection.Version == "" {
		ts.root = ts.loadRoot()
	}

	href, ok := ts.root.Collection.Links.findRelLink("teams")
	if !ok {
		return Teams{}, false
	}

	href = fmt.Sprintf("%s/search?division_id=%d", href, ts.configuration.Division)

	tr, ok := ts.makeRequest(href)
	if !ok {
		return Teams{}, false
	}

	var teams Teams
	for _, item := range tr.Collection.Items {
		if t, ok := ts.team(item); ok {
			log.Debugf("Successfully built team: %v", t)
			teams = append(teams, t)
		} else {
			log.Warnf("Failed to build team: %v", item)
		}
	}

	return teams, true
}

func (ts TeamSnap) team(team relHrefData) (Team, bool) {
	var t Team
	var results nameValueResults
	var ok bool

	// Look up the team information
	var division = string(ts.configuration.Division)
	if results, ok = team.Data.findValues("name", "id", "is_archived_season", "is_retired", "division_name"); !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warn("Unable to find specified values in return set")
		return t, false
	}

	// Determine if we should process this team
	if results["is_archived_season"] == "true" || results["is_retired"] == "true" {
		log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Skipping team because it's either archived: %s or retired: %s", results["is_archived_season"], results["is_retired"])
		return t, false
	}

	log.WithFields(log.Fields{"package": "teamsnap"}).Infof("Processing Team: %s", results["name"])

	// Grab some basic team information
	t.Name = results["name"]
	t.Level = results["division_name"]
	t.ID = generateHash(string(division), results["id"])

	// Get all of the team details
	if ok = ts.bulkLoadTeam(results["id"], &t); !ok {
		return t, false
	}

	return t, true
}

func (ts TeamSnap) bulkLoadTeam(id string, team *Team) bool {
	var href string
	var ok bool
	if href, ok = ts.root.Collection.Queries.findRelLink("bulk_load"); !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warn("Unable to find bulk_load")
		return false
	}

	// build bulk url
	var s []string
	for k := range categories {
		s = append(s, k)
	}
	bulkURL := fmt.Sprintf("%s?team_id=%s&types=%s", href, id, url.QueryEscape(strings.Join(s, ",")))

	tr, ok := ts.makeRequest(bulkURL)
	if !ok {
		return false
	}

	// Process the opponents, locations and division locations since events depends on them
	for _, item := range tr.Collection.Items {
		itemtype, _ := getItemType(item)
		category, _ := categories[itemtype]

		results, _ := loadFields(category, item)
		switch itemtype {
		case "opponent":
			id, name := processOpponent(results)
			if ts.opponents == nil {
				ts.opponents = make(map[string]string)
			}
			ts.opponents[id] = name
		case "location":
			id, location := processLocation(results)
			if ts.locations == nil {
				ts.locations = make(map[string]TeamEventLocation)
			}
			ts.locations[id] = location
		case "division_location":
			id, location := processLocation(results)
			if ts.divisionLocations == nil {
				ts.divisionLocations = make(map[string]TeamEventLocation)
			}
			ts.divisionLocations[id] = location
		case "team_photo":
			processTeamPhoto(item, team)
		}
	}

	// Process the results
	for _, item := range tr.Collection.Items {
		if ok = ts.processItem(item, team); !ok {
			return false
		}
	}

	// TODO: Sort items
	// TODO: Sort members
	// TODO: Sort events

	log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Team Object: %v", team)

	return true

}

func getItemType(item relHrefData) (string, bool) {
	var results nameValueResults
	var ok bool
	if results, ok = item.Data.findValues("type"); !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unable to load type for item: %v", item)
		return "", false
	}

	return results["type"], true

}

func loadFields(category itemType, item relHrefData) (nameValueResults, bool) {

	names := category.itemNames
	if names != nil {
		// Load the needed fields based on the type
		results, ok := item.Data.findValues(names...)
		if !ok {
			log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unable to load names: %v for item: %v", names, item)
		}
		return results, ok

	}

	return nameValueResults{}, false
}

func (ts TeamSnap) processItem(item relHrefData, team *Team) bool {

	itemtype, ok := getItemType(item)
	if !ok {
		return false
	}

	// See if we care about this category
	category, ok := categories[itemtype]
	if !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Encountered unknown item type: %s for item: %v", itemtype, item)
		return false
	}

	// Check if we care about this item
	if category.itemProcessor == nil {
		return true
	}

	// Load the values
	results, ok := loadFields(category, item)
	if !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Failed to load requested fields: %v from: %v", category, item)
		return false
	}

	// process the data
	if ok = category.itemProcessor(results, ts, team); !ok {
		log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Failed to process item: %v", item)
		return false
	}

	return true
}

func processTeamPreferences(data nameValueResults, _ TeamSnap, team *Team) bool {

	if strings.EqualFold(data["gender"], "Men") {
		team.Gender = "Boys"
	} else {
		team.Gender = "Girls"
	}

	return true
}

func processLocation(results nameValueResults) (string, TeamEventLocation) {

	return results["id"], TeamEventLocation{
		Name:    results["name"],
		Address: results["address"],
	}
}

func processOpponent(results nameValueResults) (string, string) {
	return results["id"], results["name"]
}

func processTeamPhoto(data relHrefData, team *Team) bool {

	href, _ := data.Links.findRelLink("image_url")
	team.PhotoURL = href

	log.Debugf("Team photo: %s", team.PhotoURL)

	return true
}
