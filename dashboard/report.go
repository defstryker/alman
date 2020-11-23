package dashboard

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	. "github.com/logrusorgru/aurora"
)

const (
	spcCount = "/fapi/incidents/search?count=750&gid=%s&&happenedAfter=%s&happenedBefore=%s&"
	incCount = "/fapi/incidents/search?count=20&gid=%s&&happenedAfter=%s&happenedBefore=%s&"
	benCount = "/fapi/incidents/search?count=20&gid=%s&&happenedAfter=%s&happenedBefore=%s&&fps=fp&"
	malCount = "/fapi/incidents/search?count=20&gid=%s&&happenedAfter=%s&happenedBefore=%s&&fps=tp&"

	topApps = "/fapi/stats/top-active-apps?days=30&gid="
	topNet  = "/fapi/stats/top-active-apps?days=30&eventType=8&gid="
	topUns  = "/fapi/stats/top-active-apps?days=30&eventType=2&signed=false&gid="

	endOS = "/fapi/stats/endpoints-os?gid="

	// url    = "/fapi/incidents/search?count=10&gid=%s&&happenedAfter=%s&happenedBefore=%s&"
	// urlmal = "/fapi/incidents/search?count=10&gid=%s&&happenedAfter=%s&happenedBefore=%s&&fps=tp&"
)

type TopCount []struct {
	Filename string `json:"filename"`
	Count    int    `json:"count"`
	ByType   []struct {
		Type  int `json:"type"`
		Count int `json:"count"`
	} `json:"byType"`
}

type endpointOS []struct {
	Os      string  `json:"os"`
	Count   int     `json:"count"`
	Portion float64 `json:"portion"`
}

type counter struct {
	Inc int
	Mal int
}

func (s *Server) getAlerts(month int, gid string, url string, breakdown bool) *Alerts {
	startDate := time.Date(2020, time.Month(month), 1, 0, 0, 0, 0, time.FixedZone("Asia/Singapore", 0))
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	resp, err := s.Session.Get(fmt.Sprintf(s.URL+url, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")), nil)
	fmt.Println(fmt.Sprintf(s.URL+url, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")))
	// resp, err := s.Session.Get("https://s1.managedpdr.csintelligence.asia/fapi/incidents/search?count=20&gid=547758193237819399&&happenedAfter=2020-10-01T18:30:00.000Z&happenedBefore=2020-10-31T18:30:00.000Z&", nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("1. Failed to get Incidents: " + s.Name)
	}
	var alerts = Alerts{}
	// fmt.Println(resp.StatusCode)
	resp.JSON(&alerts)
	return &alerts
}

func (s *Server) GetIncCount(month int, gid string) int {
	alerts := s.getAlerts(month, gid, incCount, true)
	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, gid))
	if err != nil {
		log.Fatalln(err)
	}
	f.NewSheet("Total Incidents")
	f.SetSheetRow("Total Incidents", "A1", &[]interface{}{alerts.RemainingItems + len(alerts.Data)})

	startDate := time.Date(2020, time.Month(month), 1, 0, 0, 0, 0, time.FixedZone("Asia/Singapore", 0))
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	bd := map[int]int{}

	for i := startDate.Day(); i <= endDate.Day(); i++ {
		st := time.Date(2020, time.Month(month), i, 0, 0, 0, 0, time.FixedZone("Asia/Singapore", 0))
		et := st.AddDate(0, 0, 1).Add(-time.Second)

		resp, err := s.Session.Get(fmt.Sprintf(s.URL+incCount, gid, st.Format("2006-01-02T15:04:05.000Z"), et.Format("2006-01-02T15:04:05.000Z")), nil)
		if err != nil {
			log.Println(err)
			log.Fatalln("1. Failed to get Incidents: " + s.Name)
		}
		var temp Alerts
		if e := resp.JSON(&temp); e != nil {
			log.Println("getAlerts breakdown")
			log.Fatalln(e)
		}
		bd[i] = temp.RemainingItems + len(temp.Data)
	}

	f.NewSheet("Incident Breakdown")

	for idx, val := range bd {
		f.SetSheetRow("Incident Breakdown", fmt.Sprintf("A%d", idx), &[]interface{}{idx, val})
	}

	if err := f.Save(); err != nil {
		println(err.Error())
	}

	return len(alerts.Data)
}

func (s *Server) GetBenignCount(month int, gid string) int {
	alerts := s.getAlerts(month, gid, benCount, false)

	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, gid))
	if err != nil {
		log.Fatalln(err)
	}
	f.NewSheet("Total Benign Count")

	f.SetSheetRow("Total Benign Count", "A1", &[]interface{}{"Total benign events", alerts.RemainingItems + len(alerts.Data)})

	if err := f.Save(); err != nil {
		println(err.Error())
	}

	return alerts.RemainingItems + len(alerts.Data)
}

func (s *Server) GetMalCount(month int, gid string) int {
	alerts := s.getAlerts(month, gid, malCount, false)

	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, gid))
	if err != nil {
		log.Fatalln(err)
	}
	f.NewSheet("Total Malicious Count")

	f.SetSheetRow("Total Malicious Count", "A1", &[]interface{}{"Total malicious events", alerts.RemainingItems + len(alerts.Data)})

	if err := f.Save(); err != nil {
		println(err.Error())
	}

	return alerts.RemainingItems + len(alerts.Data)
}

func (s *Server) GetTopCount(month int, gid string) {
	tcA := TopCount{}
	respA, err := s.Session.Get(s.URL+topApps+gid, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("1. Failed to get Incidents: " + s.Name)
	}
	e := respA.JSON(&tcA)
	if e != nil {
		fmt.Println(respA.String())
		log.Fatalln(e)
	}

	tcN := TopCount{}
	respA, err = s.Session.Get(s.URL+topNet+gid, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("1. Failed to get Incidents: " + s.Name)
	}
	respA.JSON(&tcN)

	tcU := TopCount{}
	respA, err = s.Session.Get(s.URL+topUns+gid, nil)
	if err != nil {
		log.Println(err)
		log.Fatalln("1. Failed to get Incidents: " + s.Name)
	}
	respA.JSON(&tcU)

	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, gid))
	if err != nil {
		log.Fatalln(err)
	}
	f.NewSheet("TopCount-Apps")
	f.NewSheet("TopCount-Network")
	f.NewSheet("TopCount-Unsigned")

	for idx, val := range tcA {
		f.SetSheetRow("TopCount-Apps", fmt.Sprintf("A%d", idx+1), &[]interface{}{val.Filename, val.Count})
	}

	for idx, val := range tcN {
		f.SetSheetRow("TopCount-Network", fmt.Sprintf("A%d", idx+1), &[]interface{}{val.Filename, val.Count})
	}

	for idx, val := range tcU {
		f.SetSheetRow("TopCount-Unsigned", fmt.Sprintf("A%d", idx+1), &[]interface{}{val.Filename, val.Count})
	}

	if err := f.Save(); err != nil {
		println(err.Error())
	}
}

func (s *Server) PerSiteMM(month int) {

	startDate := time.Date(2020, time.Month(month), 1, 0, 0, 0, 0, time.FixedZone("Asia/Singapore", 0))
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	sgd := map[string]string{}
	sgd["mmh"] = "548297088920715271"
	sgd["hpam"] = "548300174179434503"
	sgd["jw"] = "548299877956714503"
	sgd["pj"] = "548297435672215559"
	sgd["pj_aruba"] = "548297584599367687"
	sgd["uncat"] = "552126169328123911"
	sgd["servers"] = "562607517121642503"

	intd := map[string]string{}
	intd["mmbr"] = "551237458365251591"
	intd["mmbrsvr"] = "562943553743880199"
	intd["mmca"] = "551237354174545927"
	intd["mmcasvr"] = "562943620840161287"
	intd["mmcd"] = "551239044944625671"
	intd["mmcdsvr"] = "562943646609965063"
	intd["mmcq"] = "551237849781895175"
	intd["mmcqsvr"] = "562943668411957255"
	intd["mmcqcbz"] = "551238260559446023"
	intd["mmcqcbzsvr"] = "562943699617579015"
	intd["mmcz"] = "551237295714336775"
	intd["mmczsvr"] = "562943733453029383"
	intd["mmin"] = "551237133545766919"
	intd["mminsvr"] = "562943760296574983"
	intd["mmjz"] = "551237409329643527"
	intd["mmjzsvr"] = "562943806354227207"
	intd["mmks"] = "551237647759048711"
	intd["mmkssvr"] = "562943831587160071"
	intd["mmsh"] = "551237530486308871"
	intd["mmshsvr"] = "562943860884373511"
	intd["mmsz"] = "551238099674333191"
	intd["mmszsvr"] = "562943896192024583"
	intd["mmt"] = "561090775694180359"
	intd["mmtsvr"] = "562943926240018439"
	intd["mmxm"] = "551237800029061127"
	intd["mmxmsvr"] = "562943954723536903"
	intd["mst"] = "551237606650675207"
	intd["mstsvr"] = "562943986117902343"

	keys := make([]string, 0)
	for k := range intd {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, "547758193237819399"))
	if err != nil {
		log.Fatalln(err)
	}

	f.NewSheet("SG-Endpoints")
	f.NewSheet("Intl-Endpoints")
	var wg sync.WaitGroup

	sgdResults := map[string]counter{}
	intdResults := map[string]counter{}

	// sgd
	for name, gid := range sgd {
		wg.Add(1)

		go func(n string, gid string) {
			defer wg.Done()
			fmt.Println(n, fmt.Sprintf(s.URL+incCount, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")))
			resp, err := s.Session.Get(fmt.Sprintf(s.URL+incCount, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")), nil)
			if err != nil {
				log.Println(err)
				log.Fatalln("1. Failed to get Incidents: " + s.Name)
			}
			var als Alerts
			err = resp.JSON(&als)
			if err != nil {
				log.Println(n)
				log.Println(resp.StatusCode, resp.String())
				log.Fatalln(err)
			}

			resp, err = s.Session.Get(fmt.Sprintf(s.URL+malCount, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")), nil)
			if err != nil {
				log.Println(err)
				log.Fatalln("1. Failed to get Incidents: " + s.Name)
			}
			var als2 Alerts
			err = resp.JSON(&als2)
			if err != nil {
				log.Println(n)
				log.Println(resp.StatusCode, resp.String())
				log.Fatalln(err)
			}

			c := counter{
				Inc: len(als.Data) + als.RemainingItems,
				Mal: len(als2.Data) + als2.RemainingItems,
			}
			fmt.Println(n, c)

			sgdResults[n] = c
		}(name, gid)
	}

	// intd

	for name, gid := range intd {
		wg.Add(1)

		go func(n string, gid string) {
			defer wg.Done()
			resp, err := s.Session.Get(fmt.Sprintf(s.URL+incCount, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")), nil)
			if err != nil {
				log.Println(err)
				log.Fatalln("1. Failed to get Incidents: " + s.Name)
			}
			var als Alerts
			err = resp.JSON(&als)
			if err != nil {
				log.Println(n)
				log.Println(resp.StatusCode, resp.String())
				log.Fatalln(err)
			}

			resp, err = s.Session.Get(fmt.Sprintf(s.URL+malCount, gid, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z")), nil)
			if err != nil {
				log.Println(err)
				log.Fatalln("1. Failed to get Incidents: " + s.Name)
			}
			var als2 Alerts
			err = resp.JSON(&als2)
			if err != nil {
				log.Println(n)
				log.Println(resp.StatusCode, resp.String())
				log.Fatalln(err)
			}

			c := counter{
				Inc: len(als.Data) + als.RemainingItems,
				Mal: len(als2.Data) + als2.RemainingItems,
			}
			fmt.Println(n, c)

			intdResults[n] = c
		}(name, gid)
	}

	wg.Wait()

	counter := 1

	f.SetSheetRow("SG-Endpoints", fmt.Sprintf("A%d", counter), &[]interface{}{"Name", "Total Incidents", "Malicious Incidents"})
	counter++
	for name, val := range sgdResults {
		fmt.Println(name, val)
		f.SetSheetRow("SG-Endpoints", fmt.Sprintf("A%d", counter), &[]interface{}{name, val.Inc, val.Mal})
		counter++
	}

	counter = 1

	f.SetSheetRow("Intl-Endpoints", fmt.Sprintf("A%d", counter), &[]interface{}{"Name", "Total Incidents", "Malicious Incidents", "Server Total Incidents", "Server Malicious Incidents"})
	counter++
	for name, val := range intdResults {
		if strings.HasSuffix(name, "svr") {
			continue
		}
		svr := intdResults[name+"svr"]
		fmt.Println(name, val, name+"svr", svr)
		f.SetSheetRow("Intl-Endpoints", fmt.Sprintf("A%d", counter), &[]interface{}{name, val.Inc, val.Mal, svr.Inc, svr.Mal})
		counter++
	}

	if err := f.Save(); err != nil {
		println(err.Error())
	}
}

func (s *Server) GetEndpointOS(month int, gid string) {
	resp, err := s.Session.Get(s.URL+endOS+gid, nil)
	if err != nil {
		log.Fatalln(err)
	}

	eos := endpointOS{}
	resp.JSON(&eos)

	f, err := excelize.OpenFile(fmt.Sprintf("%d-%s.xlsx", month, gid))
	if err != nil {
		log.Fatalln(err)
	}
	f.NewSheet("Endpoints-OS")

	counter := 1
	f.SetSheetRow("Endpoints-OS", "A1", &[]interface{}{"OS", "Count", "Portion"})
	counter++

	for idx, val := range eos {
		f.SetSheetRow("Endpoints-OS", fmt.Sprintf("A%d", idx+counter), &[]interface{}{val.Os, val.Count, val.Portion})
	}

	if err := f.Save(); err != nil {
		println(err.Error())
	}
}

func (s *Server) GenReport(month int, gid string) {

	fn := fmt.Sprintf("%d-%s.xlsx", month, gid)

	f := excelize.NewFile()
	if err := f.SaveAs(fn); err != nil {
		println(err.Error())
	}
	fmt.Println(Magenta("Getting total incident count: ").String())
	s.GetIncCount(month, gid)
	// fmt.Println(Magenta("Total: ").String(), totalIncCount)
	fmt.Println(Magenta("Getting benign incident count: ").String())
	s.GetBenignCount(month, gid)
	// fmt.Println(Magenta("Benign: ").String(), benignCount)
	fmt.Println(Magenta("Getting malicious incident count: ").String())
	s.GetMalCount(month, gid)
	// fmt.Println(Magenta("Malicious: ").String(), maliciousCount)
	if gid == "547758193237819399" {
		fmt.Println(Magenta("Getting per-site breakdown for MM: ").String())
		s.PerSiteMM(month)
	}
	fmt.Println(Magenta("Getting top apps: ").String())
	s.GetTopCount(month, gid)
	fmt.Println(Magenta("Getting Endpoints OS: ").String())
	s.GetEndpointOS(month, gid)
}
