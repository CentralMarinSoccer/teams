package teamsnap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

func processMember(data nameValueResults, _ TeamSnap, team *Team) bool {
	log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Processing Member: %v", data)

	// We only care about position for non-players
	isPlayer := data["is_non_player"] == "false"
	position := data["position"]
	if isPlayer {
		position = ""
	}

	// update the team year
	teamYear(data["birthday"], &team.Year)

	// Add the member
	team.Members = append(team.Members, TeamMember{
		Name:       name(data["first_name"], data["last_name"], isPlayer),
		IsPlayer:   isPlayer,
		Position:   position,
	})

	return true
}

func name(first string, last string, isPlayer bool) string {

	// Only show initial for last name for players
	if isPlayer && len(last) > 0 {
		// Just show the first initial of the last name
		last = last[0:1]
	}

	return fmt.Sprintf("%s %s", first, last)
}

// Use birthday to determine the team's player year, e.g. 2008
func teamYear(birthday string, year *int) {

	var y int
	if len(birthday) > 4 {
		y, _ = strconv.Atoi(birthday[:4])
	}

	if y > 1990 && (y < *year || *year == 0) {
		*year = y
	}
}
