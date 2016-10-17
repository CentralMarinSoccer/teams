package geocode

import (
	"net/http"
	"net/url"
	"strings"
	"log"
	"encoding/json"
	"github.com/centralmarinsoccer/teams/cache"
)

// Geocoder holds the Google key and a cache of address string to geocoded information
type Geocoder struct {
	googleAPIKey string
	cache map[string] Address
}

// Address contains the returned address and the latitude and longitude
type Address struct {
	FormattedAddress string
	Lat float64
	Lng float64
}

// Types necessary to process google API JSON
type (
	// Response from Google
	Response struct {
		Status  string   `json:"status"`
		Results []Result `json:"results"`
	}

	// Result of the call
	Result struct {
		Types             []string           `json:"types"`
		FormattedAddress  string             `json:"formatted_address"`
		AddressComponents []AddressComponent `json:"address_components"`
		Geometry          GeometryData       `json:"geometry"`
	}

	// AddressComponent contains the address information
	AddressComponent struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	}

	// GeometryData contains geometry data
	GeometryData struct {
		Location     LatLng `json:"location"`
		LocationType string `json:"location_type"`
		Viewport     struct {
				     Southwest LatLng `json:"southwest"`
				     Northeast LatLng `json:"northeast"`
			     } `json:"viewport"`
		Bounds struct {
				     Southwest LatLng `json:"southwest"`
				     Northeast LatLng `json:"northeast"`
			     } `json:"bounds"`
	}

	// LatLng contains Latitude and Longitude values
	LatLng struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
)

const defaultFilename = "geolocation.json"

// New creates a new Geocoder and loads the cache if available
func New(key string) (Geocoder) {

	var geocoder = Geocoder{
		googleAPIKey: key,
		cache: make(map[string] Address),
	}

	cache.Load(defaultFilename, geocoder.cache)

	return geocoder
}

// SaveCache saves the cache back to disk
func (g Geocoder) SaveCache() {
	cache.Save(defaultFilename, g.cache)
}

// Lookup geocodes the specified address
func (g Geocoder) Lookup(address string) (*Address) {

	if address == "" {
		return nil
	}

	var obj Address
	var ok bool
	if obj, ok = g.cache[address]; ok {
		return &obj
	}

	resp, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?sensor=false&key=" + g.googleAPIKey + "&address=" + url.QueryEscape(strings.TrimSpace(address)))
	if err != nil {
		log.Println("Unable to contact Google. Error: " + err.Error())
		return nil
	}
	defer resp.Body.Close()

	var response = new(Response)
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Println("Unable to parse Google API response. Error: " + err.Error())
		return nil
	}

	if response.Status != "OK" {
		log.Printf("Geocoder service error: %s\n", response.Status)
		return nil
	}

	obj = Address{
		Lat: response.Results[0].Geometry.Location.Lat,
		Lng: response.Results[0].Geometry.Location.Lng,
		FormattedAddress: response.Results[0].FormattedAddress,
	}

	// Update our cache
	g.cache[address] = obj

	return &obj
}
