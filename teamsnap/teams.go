package teamsnap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

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
			teams = append(teams, t)
		}
	}

	return teams, true
}

func (ts TeamSnap) team(team relHrefData) (Team, bool) {
	var t Team

	var division = string(ts.configuration.Division)
	if results, ok := team.Data.findValues("name", "id", "is_archived_season", "is_retired", "division_name"); ok {
		if results["is_archived_season"] == "true" || results["is_retired"] == "true" {
			return t, false
		}

		log.WithFields(log.Fields{"package":"teamsnap"}).Infof("Processing Team: %s", results["name"])

		t.Name = results["name"]
		t.Level = results["division_name"]
		t.ID = generateHash(string(division), results["id"])

		ts.teamPreferences(team.Links, &t)

		t.Members, t.Year = ts.members(team.Links)
		t.Events = ts.events(team.Links)
	}

	return t, true
}
