package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/defstryker/alman/dashboard"
	"github.com/spf13/cobra"

	"github.com/gen2brain/beeep"
	. "github.com/logrusorgru/aurora"
)

var rootCmd = &cobra.Command{
	Use:   "alman",
	Short: "Alert manager",
	Long:  `Alman is an alert manager for a couple of dashboards and other services`,
}

var report = &cobra.Command{
	Use:   "report",
	Short: "Report Generation",
	Long:  `Data for report generation `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Yellow("Starting report gen...").String())
	},
}

func main() {
	// Read config file
	var cfg Config
	cfg.Read()

	rootCmd.AddCommand(report)

	// Setup cobra for getting outdated endpoints
	rootCmd.PersistentFlags().BoolP("outdated", "o", false, "When set, generates the outdated endpoints list")
	// report settings
	// report.PersistentFlags().StringP("gid", "g", "", "GID of the client")
	// report.PersistentFlags().IntP("month", "m", int(time.Now().Month()-1), "Month to get results for")
	rootCmd.PersistentFlags().StringP("gid", "g", "", "GID of the client")
	rootCmd.PersistentFlags().IntP("month", "m", int(time.Now().Month()-1), "Month to get results for")

	if cobraErr := rootCmd.Execute(); cobraErr != nil {
		panic(cobraErr)
	}

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
		// get outdated clients
		// wgAuth.Add(1)
		// go servers[server.Name].GetOutdated(&wgAuth)
	}

	wgAuth.Wait()

	// report stuff here
	// Solas: 527345349870747655
	// MM : 547758193237819399
	month, e := rootCmd.Flags().GetInt("month")
	if e != nil {
		log.Fatalln("no month")
	}
	gid, f := rootCmd.Flags().GetString("gid")
	if f != nil {
		log.Fatalln("no gid")
	}

	servers["S1"].GenReport(month, gid)
	// servers["S1"].PerSiteMM(10)
	// servers["S1"].GetTopCount("547758193237819399")

	os.Exit(1)

	// servers["S1"].GenReport(10, "547758193237819399")
	// fmt.Println(Magenta("Total: ").String(), totalIncCount)
	// benignCount := servers["S1"].GetBenignCount(10, "547758193237819399")
	// fmt.Println(Magenta("Benign: ").String(), benignCount)
	// maliciousCount := servers["S1"].GetMalCount(10, "547758193237819399")
	// fmt.Println(Magenta("Malicious: ").String(), maliciousCount)

	os.Exit(0)

	// report ends here...

	// get outdated list
	outdated, oe := rootCmd.Flags().GetBool("outdated")
	if oe != nil {
		panic(oe)
	}
	if outdated {
		var w sync.WaitGroup
		for _, s := range servers {
			w.Add(1)
			s.GetOutdated(&w)
		}
		w.Wait()
		os.Exit(0)
	}

	// waitgroup to gather the goroutines
	var wg sync.WaitGroup

	// main event loop
	for {
		// dashboards
		for _, srv := range servers {
			wg.Add(1)
			go srv.GetAlerts(&wg)
		}

		// // gmail
		// gmail := gmail.Payload{
		// 	Wg:           &wg,
		// 	ClientID:     cfg.Gmail.ClientID,
		// 	ClientSecret: cfg.Gmail.ClientSecret,
		// 	RefreshToken: cfg.Gmail.RefreshToken,
		// 	User:         cfg.Gmail.ID,
		// }

		// wg.Add(1)
		// // emails := gmail.GetUnread()

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
				stat += Cyan(fmt.Sprintf("[%s] :: %-50s: Alerts open...\n", time.Now().Format(time.RFC1123), srv.Name)).String()
			}
		}

		// gmail
		// stat += Magenta("\n\nGMail\n=====\n").String()
		// if emails == 0 {
		// 	stat += "No unread emails"
		// } else {
		// 	stat += Cyan(fmt.Sprintf("%d unread emails...\n", emails)).String()
		// }

		fmt.Print(stat)

		time.Sleep(time.Second * 20)
	}
}
