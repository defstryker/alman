package dashboard

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/xlzd/gotp"

	"github.com/levigross/grequests"
)

const (
	loginURL      string = "/login"
	verifyURL     string = "/fapi/auth/verify"
	passwordURL   string = "/fapi/auth/v2/step/password"
	challengeURL  string = "/fapi/auth/v2/step/challenge"
	openIncidents string = "/fapi/incidents/search?count=500&open=true&state=0&state=1&open=true&state=2&"
	endpoints     string = "/fapi/endpoints?sortBy=agentVersion:asc&count=500&"

	None Status = iota
	New
	Existing
)

func writeGob(filePath string, ro *grequests.RequestOptions) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(ro)
		return nil
	}
	return err
}

// Authenticate to a server, Create a new Session and return it
func (s *Server) Authenticate(wg *sync.WaitGroup) {
	defer wg.Done()
	s.Session = grequests.NewSession(nil)
	resp, err := s.Session.Get(s.URL+loginURL, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to get csrf: " + s.Name)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader((resp.String())))
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to parse html: " + s.Name)
	}

	csrfToken := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "csrf-token" {
			csrfToken, _ = s.Attr("content")
			return
		}
	})

	resp, err = s.Session.Get(s.URL+verifyURL, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to verify: " + s.Name)
	}

	s.Headers = &map[string]string{
		"csrf-token":   csrfToken,
		"Connection":   "keep-alive",
		"Content-Type": "application/json",
		"Accept":       "application/x-www-form-urlencoded; charset=utf-8",
	}

	resp, err = s.Session.Post(s.URL+passwordURL, &grequests.RequestOptions{
		JSON: map[string]string{
			"username": s.Username,
			"password": s.Password,
		},
		// Cookies: cookies.Cookies,
		Headers: *s.Headers,
	})
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to login: " + s.Name)
	}

	// pretty.Println(resp.String())
	if !s.IsThirdParty {
		// log.Println("TOTP logic here...")
		totp := gotp.NewDefaultTOTP(s.Seed)
		// log.Println(totp.Now())
		resp, err = s.Session.Post(s.URL+challengeURL, &grequests.RequestOptions{
			JSON: map[string]string{
				"username":  s.Username,
				"password":  s.Password,
				"challenge": totp.Now(),
			},
			// Cookies: cookies.Cookies,
			Headers: *s.Headers,
		})
		if err != nil {
			log.Println(err)
			log.Fatalln("Failed to login: " + s.Name)
		}
		// pretty.Println(resp.String())
	}
}

// GetAlerts gets data from each dashboard
func (s *Server) GetAlerts(wg *sync.WaitGroup) {
	defer wg.Done()

	// ensuring new alerts flag is always none before we get the alerts
	s.NewAlerts = None

	resp, err := s.Session.Get(s.URL+openIncidents, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to get data from: " + s.Name)
	}
	var alerts Alerts
	err = resp.JSON(&alerts)
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to unmarshal json in: " + s.Name)
	}

	// pretty.Println(s.Name, len(alerts.Data), len(s.Events))

	/*
		- If there are no new alerts, clear events.
		- If len of events is greater than 0,
			then send a struct to show alerts goroutine.
		- If there are new alert IDs, then add flag to send up
			a notification.
	*/
	if len(alerts.Data) == 0 {
		s.Events = nil
		s.NewAlerts = None
		return
	}

	if len(s.Events) == 0 {
		for _, data := range alerts.Data {
			s.NewAlerts = New
			s.Events = append(s.Events, data.ID)
		}
		return
	}
	s.NewAlerts = Existing

	// check if all alerts are present in current set
	for _, data := range alerts.Data {
		if !isInEvents(data.ID, &s.Events) {
			s.NewAlerts = New
			s.Events = append(s.Events, data.ID)
		}
	}
}

func isInEvents(id string, events *[]string) bool {
	for _, val := range *events {
		if val == id {
			return true
		}
	}
	return false
}

// GetOutdated clients
func (s *Server) GetOutdated(wg *sync.WaitGroup) {
	log.Println("Getting outdated from: ", s.Name)
	defer wg.Done()
	// open file
	f, err := os.Create("./outputs/" + s.Name + ".txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ge GetEndpoints

	resp, err := s.Session.Get(
		s.URL+endpoints,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = resp.JSON(&ge)
	if err != nil {
		log.Fatalln(err)
	}

	var outdated string = ""
	for _, e := range ge.Data {
		if isOutdated(e.AgentVersion) {
			outdated += fmt.Sprintf("%s, %s\n", e.Name, e.AgentVersion)
		}
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(outdated)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func isOutdated(ver string) bool {
	outdated := []string{"0.50.0", "2.0.4", "2.0.5", "2.0.14"}
	for _, out := range outdated {
		if out == ver || ver == "" {
			return true
		}
	}
	return false
}
