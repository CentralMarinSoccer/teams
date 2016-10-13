package teamsnap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"
)

func (r TeamSnap) makeRequest(url string) (teamSnapResult, bool) {
	auth := fmt.Sprintf("Bearer %s", r.configuration.Token)

	// Make sure we're pointing to the correct server (allows for mocking server for testing"
	url = strings.Replace(url, defaultServer, r.configuration.TeamSnapServer, -1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create new request for url %s\n", url)
		return teamSnapResult{}, false
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Requesting TeamSnap URL: %s\n", url)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return teamSnapResult{}, false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read complete response body")
		return teamSnapResult{}, false
	}

	if r.configuration.DumpJSON {
		newpath := "./dump"
		// Make sure the dump folder exists
		os.MkdirAll(newpath, os.ModePerm)

		tmp := strings.Replace(url, r.configuration.TeamSnapServer, "", -1)
		filename := strings.Replace(tmp, "/", "_", -1)
		fullpath := fmt.Sprintf("%s/%s.json", newpath, filename)
		ioutil.WriteFile(fullpath, body, os.ModePerm)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Request failed. Code: %d Message: %s\n", res.StatusCode, string(body[:]))
		return teamSnapResult{}, false
	}

	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var tr teamSnapResult
	if err := d.Decode(&tr); err != nil {
		log.Printf("TeamSnap JSON Root - Could not parse: %v\n", err)
		return teamSnapResult{}, false
	}

	return tr, true
}
