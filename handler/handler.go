package handler

import (
	"encoding/json"
	"net/http"
	"github.com/centralmarinsoccer/teams/teamsnap"
	"sort"
	"strings"
	"log"
	"sync"
)

const defaultPath = "/teams/"

type teams struct {
	clubData teamsnap.ClubDataInterface
	updateData chan bool
	sync.RWMutex

	teamsJSON       []byte
	teamIDToJSONMap teamIDToJSONMap
}

type teamIDToJSONMap map[string][]byte

// New creates an HTTP handler at /teams/
func New(clubData teamsnap.ClubDataInterface, updateData chan bool) (http.HandlerFunc, error) {

	team := &teams{
		clubData: clubData,
		updateData: updateData,
	}

	if err := team.optimizeFormats(); err != nil {
		return nil, err
	}

	// Wait to be notified of updated club data
	go func() {
		for {
			<-team.updateData
			team.optimizeFormats()
			log.Println("Updated Club Data")
		}
	}()

	return team.serveHTTP, nil
}

func (t *teams) optimizeFormats() error {
	cd := t.clubData.Get()

	// Sort
	sort.Sort(cd)

	// Generate our Teams JSON
	output, err := json.Marshal(cd)
	if err != nil {
		return err
	}

	// Mutex around access
	t.Lock()
	t.teamsJSON = output

	if err := t.toMap(cd); err != nil {
		return err
	}
	t.Unlock()

	return nil
}

func (t *teams) serveHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")

	ID := strings.TrimPrefix(req.URL.Path, defaultPath)
	t.RLock()
	if ID == "" {
		t.serveTeams(resp, req)
	} else {
		t.serveTeam(ID, resp, req)
	}
	t.RUnlock()
}

func (t *teams) serveTeams(resp http.ResponseWriter, req *http.Request) {

	// Return all teams
	resp.WriteHeader(200)
	resp.Write([]byte(t.teamsJSON))
}

func (t *teams) serveTeam(ID string, resp http.ResponseWriter, req *http.Request) {

	// Return the specified team
	team := t.teamIDToJSONMap[ID]
	if len(team) != 0 {
		resp.WriteHeader(200)
		resp.Write(team)
	} else {
		resp.WriteHeader(404)
	}
}

func (t *teams) toMap(clubData teamsnap.ClubData) error {
	m := make(teamIDToJSONMap)
	var err error
	for _, team := range clubData.Teams {
		if m[team.ID], err = json.Marshal(team); err != nil {
			return err
		}
	}

	t.teamIDToJSONMap = m

	return nil
}
