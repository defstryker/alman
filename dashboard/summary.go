package dashboard

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type Group struct {
	NextPage       string `json:"nextPage"`
	Data           []data `json:"data"`
	RemainingItems int    `json:"remainingItems"`
}

type data struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	Description           string        `json:"description,omitempty"`
	Created               time.Time     `json:"created"`
	EndpointNamesPreview  []string      `json:"endpointNamesPreview"`
	EndpointCount         int           `json:"endpointCount"`
	Deleted               bool          `json:"deleted"`
	Parent                string        `json:"parent,omitempty"`
	MobileEndpointsCount  int           `json:"mobileEndpointsCount"`
	RegularEndpointsCount int           `json:"regularEndpointsCount"`
	AntivirusEnabled      bool          `json:"antivirusEnabled"`
	Users                 []interface{} `json:"users"`
	License               struct {
		Limit struct {
			MaxEndpointCount       int `json:"maxEndpointCount"`
			MaxMobileEndpointCount int `json:"maxMobileEndpointCount"`
		} `json:"limit"`
		Expiration time.Time `json:"expiration"`
	} `json:"license,omitempty"`
}

func (s *Server) GetGroups() {
	resp, err := s.Session.Get(s.URL+"/fapi/profiling/groups/list", nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to get groups: " + s.Name)
	}
	var groups = Group{}
	// fmt.Println(resp.StatusCode)
	if e := resp.JSON(&groups); e != nil {
		log.Fatalln(e)
	}

	SG := map[string]int{}
	others := map[string]int{}
	total := 0

	for _, val := range groups.Data {
		if (strings.HasPrefix(val.Name, "MM") || strings.HasPrefix(val.Name, "MST")) && !strings.HasPrefix(val.Name, "MM Push") {
			if strings.HasPrefix(val.Name, "MMSG") {
				SG[val.Name] = val.EndpointCount
			} else {
				others[val.Name] = val.EndpointCount
			}
			total += val.EndpointCount
		}
	}

	keys := make([]string, 0, len(others))
	for k := range others {
		if !strings.HasSuffix(k, "Servers") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Singapore
	fmt.Println("==================Singapore======================")
	for key, val := range SG {
		fmt.Printf("%20s | %5d\n", key, val)
	}

	fmt.Println("\n\n=================================================\n\n")

	// Others
	fmt.Println("====================Others======================")
	for _, key := range keys {
		fmt.Printf("%20s | %5d | %20s | %5d\n", key, others[key], key+" Servers", others[key+" Servers"])
	}

	fmt.Println("\n\n=================================================\n\n")
	fmt.Println("total endpoint count: ", total)
}

// func (s *Server)
