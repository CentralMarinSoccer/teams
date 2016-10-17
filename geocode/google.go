package geocode

import (
	"net/http"
	"net/url"
	"strings"
	"log"
	"encoding/json"
	"strconv"
	"github.com/centralmarinsoccer/team/cache"
)

type Geocoder struct {
	googleApiKey string
	cache map[string] Address
}

type Address struct {
	FormattedAddress string
	Lat float64
	Lng float64
}

// Types necessary to process google API JSON
type (
	Response struct {
		Status  string   `json:"status"`
		Results []Result `json:"results"`
	}

	Result struct {
		Types             []string           `json:"types"`
		FormattedAddress  string             `json:"formatted_address"`
		AddressComponents []AddressComponent `json:"address_components"`
		Geometry          GeometryData       `json:"geometry"`
	}

	AddressComponent struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	}

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

	LatLng struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
)

const defaultFilename = "geolocation.json"

func New(key string) (Geocoder) {

	var geocoder = Geocoder{
		googleApiKey: key,
		cache: make(map[string] Address),
	}

	cache.Load(defaultFilename, geocoder.cache)

	return geocoder
}

func (g Geocoder) SaveCache() {
	cache.Save(defaultFilename, g.cache)
}

func (g Geocoder) Lookup(address string) (*Address) {

	if address == "" {
		return nil
	}

	var obj Address
	var ok bool
	if obj, ok = g.cache[address]; ok {
		return &obj
	}

	resp, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?sensor=false&key=" + g.googleApiKey + "&address=" + url.QueryEscape(strings.TrimSpace(address)))
	log.Println("https://maps.googleapis.com/maps/api/geocode/json?sensor=false&key=" + g.googleApiKey + "&address=" + url.QueryEscape(strings.TrimSpace(address)))
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

	// convert floats to strings
	//lat := floatToString(response.Results[0].Geometry.Location.Lat)
	//lng := floatToString(response.Results[0].Geometry.Location.Lng)

	obj = Address{
		Lat: response.Results[0].Geometry.Location.Lat,
		Lng: response.Results[0].Geometry.Location.Lng,
		FormattedAddress: response.Results[0].FormattedAddress,
	}

	// Update our cache
	g.cache[address] = obj

	return &obj
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
