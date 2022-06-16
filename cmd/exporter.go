package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/charlie-haley/omada_exporter/pkg/omada"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var version = "development"
var (
	host                     string
	username                 string
	password                 string
	port                     string
	site                     string
	interval                 int
	timeout                  int
	insecure                 bool
	goCollectorDisabled      bool
	processCollectorDisabled bool
)

func Run() {
	app := cli.NewApp()
	app.Name = "omada_exporter"
	app.Version = version
	app.Usage = "Prometheus Exporter for TP-Link Omada Controller SDN."
	app.EnableBashCompletion = true
	app.Authors = []*cli.Author{
		{Name: "Charlie Haley", Email: "charlie-haley@users.noreply.github.com"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Destination: &host, Required: true, Name: "host", Value: "", Usage: "The hostname of the Omada Controller, including protocol.", EnvVars: []string{"OMADA_HOST"}},
		&cli.StringFlag{Destination: &username, Required: true, Name: "username", Value: "", Usage: "Username of the Omada user you'd like to use to fetch metrics.", EnvVars: []string{"OMADA_USER"}},
		&cli.StringFlag{Destination: &password, Required: true, Name: "password", Value: "", Usage: "Password for your Omada user.", EnvVars: []string{"OMADA_PASS"}},
		&cli.StringFlag{Destination: &port, Name: "port", Value: "9202", Usage: "Port on which to expose the Prometheus metrics.", EnvVars: []string{"OMADA_PORT"}},
		&cli.StringFlag{Destination: &site, Name: "site", Value: "Default", Usage: "Omada site to scrape metrics from.", EnvVars: []string{"OMADA_SITE"}},
		&cli.IntFlag{Destination: &interval, Name: "interval", Value: 5, Usage: "Interval between scrapes, in seconds.", EnvVars: []string{"OMADA_SCRAPE_INTERVAL"}},
		&cli.IntFlag{Destination: &timeout, Name: "timeout", Value: 15, Usage: "Timeout when making requests to the Omada Controller.", EnvVars: []string{"OMADA_REQUEST_TIMEOUT"}},
		&cli.BoolFlag{Destination: &insecure, Name: "insecure", Value: false, Usage: "Whether to skip verifying the SSL certificate on the controller.", EnvVars: []string{"OMADA_INSECURE"}},
		&cli.BoolFlag{Destination: &goCollectorDisabled, Name: "disable-go-collector", Value: false, Usage: "Disable Go collector metrics.", EnvVars: []string{"OMADA_DISABLE_GO_COLLECTOR"}},
		&cli.BoolFlag{Destination: &processCollectorDisabled, Name: "disable-process-collector", Value: false, Usage: "Disable process collector metrics.", EnvVars: []string{"OMADA_DISABLE_PROCESS_COLLECTOR"}},
	}
	app.Commands = []*cli.Command{
		{Name: "version", Aliases: []string{"v"}, Usage: "prints the current version.",
			Action: func(c *cli.Context) error {
				fmt.Println(version)
				os.Exit(0)
				return nil
			}},
	}
	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	if goCollectorDisabled {
		// remove Go collector
		prometheus.Unregister(prometheus.NewGoCollector())
	}
	if processCollectorDisabled {
		// remove Process collector
		prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	}

	// check if host is properly formatted
	if strings.HasSuffix(host, "/") {
		// remove trailing slash if it exists
		host = strings.TrimRight(host, "/")
	}

	client, err := api.Configure(c)
	if err != nil {
		return err
	}

	go handleScrape(client)

	log.Info(fmt.Sprintf("listening on :%s", port))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
    <head>
	<title>omada_exporter</title>
	</head>
    	<body>
			<h1>omada_exporter</h1>
			<p>
				<a href="/metrics">Metrics</a>
			</p>
    	</body>
    </html>`))
	})

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		return err
	}

	return nil
}

func handleScrape(config *api.Client) {
	// ensure we scrape before the interval
	err := omada.Scrape(config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	//nolint:gosimple
	for {
		select {
		case <-ticker.C:
			//nolint:errcheck
			go omada.Scrape(config)
		}
	}
}
