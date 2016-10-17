package main_test

//import (
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/onsi/gomega/ghttp"
//
//	"io/ioutil"
//	"net/http"
//	"os"
//	"github.com/centralmarinsoccer/teams"
//	"github.com/centralmarinsoccer/teams/handlers"
//)
//
//var _ = Describe("Team", func() {
//
//	It("responds to HTTP GET with a list of teams", func() {
//
//		filename := "teams.json"
//		server := ghttp.NewServer()
//
//		h, err := handlers.GetHandler(&handlers.Teams{
//			Path:     "/",
//			Filename: filename,
//			Data:     &teamsdata.Teams{},
//		})
//		Expect(err).NotTo(HaveOccurred())
//
//		server.AppendHandlers(h)
//		defer server.Close()
//
//		request, err := http.NewRequest("GET", server.URL(), nil)
//		Expect(err).NotTo(HaveOccurred())
//
//		res, err := http.DefaultClient.Do(request)
//		defer res.Body.Close()
//
//		Expect(err).NotTo(HaveOccurred())
//		Expect(res.StatusCode).To(Equal(200))
//
//		_, err = ioutil.ReadAll(res.Body)
//		Expect(err).NotTo(HaveOccurred())
//	})
//
//	Context("Environment Variables", func() {
//
//		BeforeEach(func() {
//			// Clear all environment variables
//			os.Setenv("DIVISION", "")
//			os.Setenv("TOKEN", "")
//			os.Setenv("PORT", "")
//			os.Setenv("DOMAIN", "")
//			os.Setenv("INTERVAL", "")
//		})
//
//		It("return false if DIVISION and TOKEN are not set", func() {
//			_, ok := main.ValidateEnvs()
//			Expect(ok).ToNot(BeTrue())
//		})
//
//		It("validates that DIVISION, INTERVAL, and PORT are integers", func() {
//			os.Setenv("DIVISION", "A1234")
//			os.Setenv("INTERVAL", "CBA")
//			os.Setenv("PORT", "098AB")
//			env, ok := main.ValidateEnvs()
//			Expect(ok).ToNot(BeTrue())
//			Expect(env.Division).To(BeEquivalentTo(-1))
//			Expect(env.RefreshInterval).To(BeEquivalentTo(-1))
//			Expect(env.URL).To(BeEquivalentTo(":-1"))
//		})
//
//		It("properly retreives DIVISION, TOKEN, INTERVAL, PORT, and DOMAIN", func() {
//
//			os.Setenv("DIVISION", "1234")
//			os.Setenv("TOKEN", "1234")
//			os.Setenv("PORT", "8000")
//			os.Setenv("DOMAIN", "example.domain.com")
//			os.Setenv("INTERVAL", "10")
//
//			env, ok := main.ValidateEnvs()
//			Expect(ok).To(BeTrue())
//			Expect(env.Division).To(BeEquivalentTo(1234))
//			Expect(env.Token).To(BeEquivalentTo("1234"))
//			Expect(env.URL).To(BeEquivalentTo("http://example.domain.com:8000"))
//
//		})
//	})
//})
