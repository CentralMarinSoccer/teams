package teamsnap_test

import (
	. "github.com/centralmarinsoccer/teams/teamsnap"

	//"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"github.com/centralmarinsoccer/teams/mocks"
	"net/http"
	"net/http/httptest"
	"log"
	"strings"
	"fmt"
)

const testToken = "ABCD"

func fixtureJson(urlPath string) ([]byte, int) {

	filename := strings.Replace(urlPath, "/", "_", -1)
	fullpath := fmt.Sprintf("../fixtures/%s.json", filename)
	fileJson, err := ioutil.ReadFile(fullpath)
	if err != nil {
		log.Printf("Failed to read file: %s. Error: %v\n", fullpath, err)
		return nil, http.StatusNotFound
	}

	return fileJson, http.StatusOK
}

var _ = Describe("TeamSnap", func() {

	var (
		cacheMfs mocks.FileSystem
		httpMfs  mocks.FileSystem

		cacheConfiguration Configuration
		httpConfiguration  Configuration

		//server1 *ghttp.Server
		server *httptest.Server
	)

	BeforeEach(func() {

		//server1 = ghttp.NewServer()
		//server1.AppendHandlers(
		//	ghttp.CombineHandlers(
		//		ghttp.VerifyRequest("GET", "/v3/"),
		//		ghttp.VerifyHeader(http.Header{
		//			"Authentication": []string{fmt.Sprintf("Bearer %s", testToken)},
		//		}),
		//		ghttp.RespondWith(http.StatusOK, `[
		//
		//		]`),
		//	),
		//)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}
			fixtureJSON, status := fixtureJson(path)
			w.WriteHeader(status)
			w.Write(fixtureJSON)
		}))

		cacheConfiguration = Configuration{Token: testToken, Division: 1234, FileSystem: &cacheMfs, TeamSnapServer: server.URL}

		httpMfs.StatCall.Returns.Error = errors.New("File not found")
		httpConfiguration = Configuration{Token: testToken, Division: 5678, FileSystem: &httpMfs, TeamSnapServer: server.URL}
	})
/*
	It("should return cached TeamSnap data", func() {

		// Load fixture data
		var err error
		cacheMfs.ReadFileCall.Returns.Data, err = ioutil.ReadFile("../fixtures/clubdata.json")

		ts, err := New(&cacheConfiguration)

		// Make sure our mock was called and we didn't get an error
		Expect(cacheMfs.StatCall.Receives.Filename).ToNot(BeEmpty(), "Mock filesystem's Stat not called")
		Expect(err).To(BeNil(), "TeamSnap New returned an error %v", err)

		// convert our fixture data to a struct
		var expected ClubData
		err = json.Unmarshal(cacheMfs.ReadFileCall.Returns.Data, &expected)
		Expect(err).To(BeNil(), "Unable to parse Fixture ClubData")

		actual := ts.Get()
		Expect(actual).Should(Equal(expected), "Fixture data doesn't match cache data")
	})
*/
	// Test HTTP Requests
	It("should return club data from mocked http calls", func() {
		_, err := New(&httpConfiguration)

		Expect(err).To(BeNil(), "TeamSnap New returned an error %v", err)
	})

	// TODO: Test bad cache file (invalid JSON)
})
