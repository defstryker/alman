package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/defstryker/alman/dashboard"
	"github.com/defstryker/alman/gmail"

	"github.com/gen2brain/beeep"
	. "github.com/logrusorgru/aurora"
)

func main() {
	// Read config file
	var cfg Config
	cfg.Read()

	// goroutines for auth
	var wgAuth sync.WaitGroup

	names := []string{}

	servers := make(map[string]*dashboard.Server)
	for _, server := range cfg.Dashboards {
		servers[server.Name] = &dashboard.Server{
			Name:         server.Name,
			URL:          server.URL,
			Username:     server.Username,
			Password:     server.Password,
			IsThirdParty: server.IsThirdParty,
			Seed:         server.Seed,
		}
		names = append(names, server.Name)

		wgAuth.Add(1)
		// run auth to get and save headers
		go servers[server.Name].Authenticate(&wgAuth)
	}

	wgAuth.Wait()

	// waitgroup to gather the goroutines
	var wg sync.WaitGroup

	// main event loop
	for {
		// dashboards
		for _, srv := range servers {
			wg.Add(1)
			go srv.GetAlerts(&wg)
		}

		// gmail
		gmail := gmail.Payload{
			Wg:           &wg,
			ClientID:     cfg.Gmail.ClientID,
			ClientSecret: cfg.Gmail.ClientSecret,
			RefreshToken: cfg.Gmail.RefreshToken,
			User:         cfg.Gmail.ID,
		}

		wg.Add(1)
		emails := gmail.GetUnread()

		wg.Wait()

		screenClear()

		var stat string

		// dashboards
		stat += Magenta("Dashboards\n==========\n").String()
		for _, sn := range names {
			srv := servers[sn]
			if srv.NewAlerts == dashboard.None {
				stat += fmt.Sprintf("[%s] :: %-20s: No new incidents\n", time.Now().Format(time.RFC1123), srv.Name)
			} else if srv.NewAlerts == dashboard.New {
				stat += Magenta(fmt.Sprintf("[%s] :: %-20s: New Alert!!!\n", time.Now().Format(time.RFC1123), srv.Name)).String()

				// beep here
				go beep()

				err := beeep.Notify(srv.Name, "New Alerts on "+srv.Name, "assets/warning.png")
				if err != nil {
					panic(err)
				}
			} else {
				stat += Cyan(fmt.Sprintf("[%s] :: %-20s: Alerts open...\n", time.Now().Format(time.RFC1123), srv.Name)).String()
			}
		}

		// gmail
		stat += Magenta("\n\nGMail\n=====\n").String()
		if emails == 0 {
			stat += "No unread emails"
		} else {
			stat += Cyan(fmt.Sprintf("%d unread emails...\n", emails)).String()
		}

		fmt.Print(stat)

		time.Sleep(time.Second * 5)
	}
}
