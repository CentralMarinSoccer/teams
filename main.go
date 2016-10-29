package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/centralmarinsoccer/teams/geocode"
	"github.com/centralmarinsoccer/teams/handler"
	"github.com/centralmarinsoccer/teams/teamsnap"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

const defaultPort = 8080
const urlPath = "/teams/"
const defaultInterval = 60

// Environment contains all of the required and optional environment variables
type Environment struct {
	Token           string
	Division        int
	URL             string
	RefreshInterval time.Duration
	GoogleAPIKey    string
}

func init() {
	loglevel := flag.Int("loglevel", 1, "Valid levels are: 1 Info, 2 Debug")
	flag.Parse()

	if *loglevel == 2 {
		log.Println("Setting Debug logging level")
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	// Make sure we have the appropriate environment variables
	var env Environment
	var ok bool
	if env, ok = ValidateEnvs(); !ok {
		log.Errorln("Missing required environment variables")
	}

	// Create a geocoder
	geocoder := geocode.New(env.GoogleAPIKey)

	ts, err := teamsnap.New(&teamsnap.Configuration{
		Division: env.Division,
		Token:    env.Token,
		Geocoder: geocoder,
	})
	if err != nil {
		log.Errorf("Failed to create a new TeamSnap. Error: %v", err)
	}

	// create a channel to update TeamSnap data
	update := make(chan bool)

	// Setup our HTTP Server
	h, err := handler.New(ts, update)
	if err != nil {
		log.Errorf("Handlers error: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/teams/", h)
	r.HandleFunc("/teams/{ID}", h)
	r.Handle("/metrics", prometheus.Handler()) // Add Metrics Handler
	r.PathPrefix("/teams/static/").Handler(http.StripPrefix("/teams/static/", http.FileServer(http.Dir("static"))))

	// Setup timer to refresh TeamSnap data
	ticker := time.NewTicker(env.RefreshInterval * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infoln("Updating team data")
				if ok := ts.Update(); ok {
					update <- true
					log.Infoln("Complete")
				} else {
					log.Warnln("Failed to update data")
				}
			}
		}
	}()

	log.Infof("Starting up server at %s%s with data refresh interval of %d for TeamSnap division %d", env.URL, urlPath, env.RefreshInterval, env.Division)
	srv := &http.Server{
		Handler: r,
		Addr:    env.URL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// ValidateEnvs checks to make sure the necessary environment variables are set
func ValidateEnvs() (Environment, bool) {

	var env Environment
	state := true

	// Required Parameters
	var ok bool
	if env.Division, ok = getEnvInt("DIVISION", 0, true); !ok {
		state = false
	}
	if env.Token, ok = getEnvString("TOKEN", "", true); !ok {
		state = false
	}

	var port int
	var domain string
	if port, ok = getEnvInt("PORT", defaultPort, false); !ok {
		state = false
	}
	if domain, ok = getEnvString("DOMAIN", "", false); !ok {
		state = false
	}

	if len(domain) != 0 {
		domain = "http://" + domain
	}
	env.URL = fmt.Sprintf("%s:%d", domain, port)

	var interval int
	if interval, ok = getEnvInt("REFRESHINTERVAL", defaultInterval, false); ok {
		env.RefreshInterval = time.Duration(interval)
	} else {
		state = false
	}

	if env.GoogleAPIKey, ok = getEnvString("GOOGLE_API_KEY", "", true); !ok {
		state = false
	}

	return env, state
}

func getEnvString(name string, defaultVal string, required bool) (string, bool) {
	val := os.Getenv(name)
	if val == "" {
		if required {
			log.Warnf("Missing required environment variable '%s'", name)
			return "", false
		}
		return defaultVal, true
	}

	return val, true
}

func getEnvInt(name string, defaultVal int, required bool) (int, bool) {

	if val, ok := getEnvString(name, "", required); ok {
		if val == "" {
			return defaultVal, true
		}

		i, err := strconv.Atoi(val)
		if err != nil {
			log.Warnf("%s's value %s is not an integer.", name, val)
			return -1, false
		}
		return i, true
	}

	return -1, false
}
