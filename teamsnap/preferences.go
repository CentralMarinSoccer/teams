package teamsnap

import "strings"

func (ts TeamSnap) teamPreferences(links relHrefDatas, t *Team) {

	// Get Team Preferences url
	if href, ok := links.findRelLink("team_preferences"); ok {

		// Get team preferences
		tr, _ := ts.makeRequest(href)

		t.PhotoURL = teamImage(tr)
		t.Gender = teamGender(tr)
	}
}

func teamImage(tr teamSnapResult) string {
	if href, ok := tr.Collection.Items[0].Links.findRelLink("team_photo"); ok {
		return href
	}

	return ""
}

func teamGender(tr teamSnapResult) string {
	if results, ok := tr.Collection.Items[0].Data.findValues("gender"); ok {
		// convert to boys / girls
		if strings.EqualFold(results["gender"], "Men") {
			return "Boys"
		} else {
			return "Girls"
		}
	}
	return ""
}
