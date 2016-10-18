package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/centralmarinsoccer/teams/handler"
	"github.com/centralmarinsoccer/teams/teamsnap"

	"net/http"
	"time"
	"github.com/centralmarinsoccer/teams/geocode"
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
	GoogleAPIKey	string
}

func main() {
	// Make sure we have the appropriate environment variables
	var env Environment
	var ok bool
	if env, ok = ValidateEnvs(); !ok {
		os.Exit(-2)
	}

	// Create a geocoder
	geocoder := geocode.New(env.GoogleAPIKey)

	ts, err := teamsnap.New(&teamsnap.Configuration{
		Division:      env.Division,
		Token:         env.Token,
		Geocoder:  geocoder,
	})
	if err != nil {
		log.Printf("Failed to create a new TeamSnap. Error: %v\n", err)
		os.Exit(-1)
	}

	// create a channel to update TeamSnap data
	update := make(chan bool)

	// Setup our HTTP Server
	mux := http.NewServeMux()
	h, err := handler.New(ts, update)

	if err != nil {
		log.Printf("Handlers error: %v\n", err)
		os.Exit(-2)
	}

	mux.Handle(urlPath, h)
	mux.Handle("/metrics", prometheus.Handler()) // Add Metrics Handler

	log.Printf("Starting up server at %s%s with data refresh interval of %d for TeamSnap division %d\n", env.URL, urlPath, env.RefreshInterval, env.Division)

	ticker := time.NewTicker(env.RefreshInterval * time.Minute)
	go func() {
		for {
			select {
			case <- ticker.C:
				log.Println("Updating data")
				if ok := ts.Update(); ok {
					update <- true
					log.Println("Complete")
				} else {
					log.Println("Failed to update data")
				}
			}
		}
	}()

	http.ListenAndServe(env.URL, mux)

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
			log.Printf("Missing required environment variable '%s'\n", name)
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
			log.Printf("%s's value %s is not an integer.\n", name, val)
			return -1, false
		}
		return i, true
	}

	return -1, false
}
