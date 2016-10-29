package teamsnap

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

func processMember(data nameValueResults, _ TeamSnap, team *Team) bool {
	log.WithFields(log.Fields{"package": "teamsnap"}).Debugf("Processing Member: %v", data)

	mt := memberType(data)

	// update the team year
	teamYear(data["birthday"], &team.Year)

	// Add the member
	team.Members = append(team.Members, TeamMember{
		Name:       name(data["first_name"], data["last_name"], mt),
		MemberType: mt,
	})

	return true
}

func memberType(results nameValueResults) string {
	if results["is_non_player"] == "false" {
		return memberTypePlayer
	}

	// Determine if this is a coach
	if strings.EqualFold(results["position"], "coach") {
		return memberTypeCoach
	}

	if results["is_owner"] == "true" {
		return memberTypeManager
	}

	return memberTypeAssistantManager
}

func name(first string, last string, mt string) string {

	// Only show initial for last name for players
	if mt == memberTypePlayer && len(last) > 0 {
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
