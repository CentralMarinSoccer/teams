package teamsnap

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/centralmarinsoccer/teams/filesystem"
	"time"
	"github.com/centralmarinsoccer/teams/geocode"
	"github.com/centralmarinsoccer/teams/cache"
)

// ClubDataInterface provides a mechanism to retrieve club data
type ClubDataInterface interface {
	Get() ClubData
}

// TeamSnap provides a mechanism to configure and cache team data stored in TeamSnap.com website
type TeamSnap struct {
	root          teamSnapResult
	locations     map[string]TeamEventLocation
	clubData      ClubData
	configuration Configuration
}

// Configuration defines how TeamSnap information should be accessed
type Configuration struct {
	Token           string
	Division        int
	Geocoder	geocode.Geocoder
	FileSystem      filesystem.LocalDiskInterface
	TeamSnapServer  string
	DumpJSON        bool
}

// ClubData is the data returned from TeamSnap
type ClubData struct {
	LastUpdated time.Time `json:"last_updated"`
	Teams       Teams     `json:"teams"`
}

// Teams is an array of Team
type Teams []Team

// Team stores the details of a Team
type Team struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Gender   string       `json:"gender"`
	Year     int          `json:"year"`
	Level    string       `json:"level"`
	PhotoURL string       `json:"image_url,omitempty"`
	Members  []TeamMember `json:"members,omitempty"`
	Events   []TeamEvent  `json:"events,omitempty"`
}

const memberTypePlayer = "player"
const memberTypeCoach = "coach"
const memberTypeManager = "manager"

// TeamMember holds the text and metadata for a team member
type TeamMember struct {
	Name       string `json:"name"`
	MemberType string `json:"type"`
}

// TeamEvent contains all of the data that makes up an event
type TeamEvent struct {
	Start     time.Time            `json:"start"`
	Opponent  string               `json:"opponent"`
	Duration  string               `json:"duration"`
	Location  TeamEventLocation    `json:"location"`
}

// TeamEventLocation contains all of the data that makes up an event location
type TeamEventLocation struct {
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Latitude float64   `json:"latitude,omitempty"`
        Longitude float64  `json:"longitude,omitempty"`
}

const defaultServer = "https://api.teamsnap.com"
const defaultFilename = "teamsnap.json"

// TODO: Notify callers when data in cache has changed. Maybe hash our data struct so we can detect changes

// New creates a TeamSnap with the specified configuration
func New(configuration *Configuration) (*TeamSnap, error) {

	// Use defaults if variable not specified
	if configuration.FileSystem == nil {
		configuration.FileSystem = filesystem.OSFS{}
	}
	if configuration.TeamSnapServer == "" {
		configuration.TeamSnapServer = defaultServer
	}

	ts := &TeamSnap{
		configuration: *configuration,
		locations: make(map[string]TeamEventLocation),
	}

	// Check if the file exists
	var dataLoaded = false
	if err := cache.Load(defaultFilename, &ts.clubData); err == nil {
		dataLoaded = true
	}

	if !dataLoaded {
		log.WithFields(log.Fields{"package":"teamsnap"}).Warnf("TeamSnap cache '%s' does not exist or failed to load. Building initial version", defaultFilename)
		dataLoaded = ts.loadTeamSnapData()
	}

	// Make sure we successfully loaded data
	if !dataLoaded {
		return &TeamSnap{}, errors.New("Failed to load TeamSnap data. Please check previous errors to determine cause.")
	}

	return ts, nil
}

func (ts *TeamSnap) loadTeamSnapData() bool {

	// Load data from TeamSnap web API
	teams, ok := ts.teams()
	if !ok {
		log.WithFields(log.Fields{"package":"teamsnap"}).Warnf("Unable to retrieve data from TeamSnap. Check previous errors")
		return false
	}

	ts.clubData.Teams = teams
	ts.clubData.LastUpdated = time.Now()

	// Save our caches
	ts.configuration.Geocoder.SaveCache()
	cache.Save(defaultFilename, ts.clubData)

	return true
}

// Get returns club data
func (ts *TeamSnap) Get() ClubData {
	return ts.clubData
}

// Update reloads the club data from TeamSnap.com
func (ts *TeamSnap) Update() bool {
	return ts.loadTeamSnapData()
}

// Implement the Sort Interface

// Len provides the number of Teams
func (cd ClubData) Len() int { return len(cd.Teams) }

// Swap swaps two teams
func (cd ClubData) Swap(i, j int) { cd.Teams[i], cd.Teams[j] = cd.Teams[j], cd.Teams[i] }

// Sort teams by gender, year, and name
func (cd ClubData) Less(i, j int) bool {
	if cd.Teams[i].Gender == cd.Teams[j].Gender {
		if cd.Teams[i].Year == cd.Teams[j].Year {
			return cd.Teams[i].Name < cd.Teams[j].Name
		}
		return cd.Teams[i].Year > cd.Teams[j].Year
	}
	return cd.Teams[i].Gender < cd.Teams[j].Gender
}
