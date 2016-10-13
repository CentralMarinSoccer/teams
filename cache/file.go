package cache

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
)

func Load(filename string, structure interface{}) error {
	if _, err := os.Stat(filename); err != nil {
		log.Printf("Cache file '%s' doesn't exist.\n", filename)
		return err
	}

	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filename); err != nil {
		log.Printf("Failed to loading data from cache file '%s'.\n", filename)
		return err
	}

	if err := json.Unmarshal(data, &structure); err != nil {
		log.Printf("Unable to parse JSON from cache file '%s'. Error: %v\n", filename, err)
		return err
	}

	return nil
}

func Save(filename string, structure interface{}) error {

	// Remove the old cache file
//	os.Remove(filename)

	// Save the data to file after successful retrieval
	resultJson, err := json.Marshal(structure)
	if err != nil {
		log.Printf("Warning: Unable to generate JSON for cache '%s'. Error: %v\n", filename, err)
		return err
	}

	if err := ioutil.WriteFile(filename, resultJson, 0644); err != nil {
		log.Printf("Warning: Unable to save cache data to file '%s'. Error: %v\n", filename, err)
		return err
	}

	return nil
}