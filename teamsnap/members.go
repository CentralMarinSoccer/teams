package teamsnap

import (
	"fmt"
	"strconv"
	"strings"
	"log"
)

func (ts TeamSnap) members(links relHrefDatas) ([]TeamMember, int) {
	year := 0
	var members []TeamMember

	if href, ok := links.findRelLink("members"); ok {
		tr, _ := ts.makeRequest(href)
		for _, m := range tr.Collection.Items {
			m1, birthday := ts.member(m)
			members = append(members, m1)

			teamYear(birthday, &year)
		}
	}

	return members, year
}

// Use birthday to determine the team's player year, e.g. 2008
func teamYear(birthday string, year *int) {

	var y int
	if len(birthday) > 4 {
		y, _ = strconv.Atoi(birthday[:4])
	}

	log.Printf("Birthday: %v, Year: %d, Current Year: %d\n", birthday, y, *year)
	if y > 1990 && (y < *year || *year == 0) {
		*year = y
		log.Println("NEW YEAR Selected")
	}
}

func (ts TeamSnap) member(m relHrefData) (TeamMember, string) {
	href := emailHref(m.Links)
	if results, ok := m.Data.findValues("first_name", "last_name", "birthday", "is_manager", "is_non_player", "is_owner"); ok {
		mt := ts.memberType(href, results)
		return TeamMember{
			Name:       name(results["first_name"], results["last_name"], mt),
			MemberType: mt,
		}, results["birthday"]
	}
	return TeamMember{}, ""
}

func (ts TeamSnap) memberType(href string, results nameValueResults) string {
	if results["is_non_player"] == "false" {
		return memberTypePlayer
	}

	// Determine if this is a coach
	if tr, ok := ts.makeRequest(href); ok {
		for _, e := range tr.Collection.Items {
			if l, ok := e.Data.findValues("label"); ok {
				if caseInsensitiveContains(l["label"], "coach") {
					return memberTypeCoach
				}
			}
		}
	}

	return memberTypeManager
}

func name(first string, last string, mt string) string {
	switch mt {
	case memberTypeCoach, memberTypeManager:
		// Show full last name, so nothing to do
	case memberTypePlayer:
		if len(last) > 0 {
			// Just show the first initial of the last name
			last = last[0:1]
		}
	}

	return fmt.Sprintf("%s %s", first, last)
}

func emailHref(links relHrefDatas) string {
	// find the emails link
	if href, ok := links.findRelLink("member_email_addresses"); ok {
		return href
	}

	return ""
}

func caseInsensitiveContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}
