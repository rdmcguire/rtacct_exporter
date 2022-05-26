package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	realmStat struct {
		Realm     string
		BytesTo   int
		PktsTo    int
		BytesFrom int
		PktsFrom  int
	}
	realmStats []*realmStat
)

// Configuration
var (
	listen   string = ":9987"
	interval string = "30s"
	delay    time.Duration
	verbose  bool
)

// Prometheus Metrics
var (
	reg      = prometheus.NewRegistry()
	metrics  = promauto.With(reg)
	promPkts = metrics.NewCounterVec(prometheus.CounterOpts{
		Name: "rtacct_pkts",
		Help: "Counter of packets from/to realm",
	}, []string{"realm", "direction"})
	promBytes = metrics.NewCounterVec(prometheus.CounterOpts{
		Name: "rtacct_bytes",
		Help: "Counter of bytes from/to realm",
	}, []string{"realm", "direction"})
	promStatsDuration = metrics.NewHistogram(prometheus.HistogramOpts{
		Name:    "rtacct_stats_duration_us",
		Help:    "Time spend gathering statistics in microseconds",
		Buckets: prometheus.LinearBuckets(500, 500, 20),
	})
)

func init() {
	// Config
	flag.StringVar(&interval, "interval", interval, "Stats update interval")
	flag.StringVar(&listen, "listen", listen, "Prometheus http listen address")
	flag.BoolVar(&verbose, "verbose", verbose, "Be Verbose")
	flag.Parse()

	// Parse Interval
	var err error
	delay, err = time.ParseDuration(interval)
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
}

func main() {
	// Start Prometheus
	go http.ListenAndServe(listen, nil)
	if verbose {
		log.Println("Started Prometheus")
	}

	// Update forever
	timer := time.NewTicker(delay)
	update() // Update before first tick
	for range timer.C {
		update()
	}
}

func update() {
	if verbose {
		log.Println("Gathering statistics from rtacct")
	}

	// Gather statistics
	t1 := time.Now()
	stats, _ := getStats()
	dur := time.Now().Sub(t1)

	// Report statistics
	promStatsDuration.Observe(float64(dur.Microseconds()))
	for _, stat := range *stats {
		promPkts.With(prometheus.Labels{"realm": stat.Realm, "direction": "from"}).Add(float64(stat.PktsFrom))
		promPkts.With(prometheus.Labels{"realm": stat.Realm, "direction": "to"}).Add(float64(stat.PktsTo))
		promBytes.With(prometheus.Labels{"realm": stat.Realm, "direction": "from"}).Add(float64(stat.BytesFrom))
		promBytes.With(prometheus.Labels{"realm": stat.Realm, "direction": "to"}).Add(float64(stat.BytesTo))
	}
}

func getStats() (*realmStats, error) {
	stats := new(realmStats)
	var err error
	out, err := exec.Command("rtacct").Output()

	if err != nil {
		log.Printf("Error running rtacct: %s", err)
		return stats, err
	}

	lines := bytes.Split(out, []byte("\n"))

	for _, line := range lines {
		if len(line) < 50 {
			continue
		} else if line[0] == ' ' {
			continue
		} else if line[11] == 'B' {
			continue
		}

		r := regexp.MustCompile(`([^\s]+)\s+(\d+[KMGT]?)\s+(\d+[KMGT]?)\s+(\d+[KMGT]?)\s+(\d+[KMGT]?)`)
		matches := r.FindAllSubmatch(line, -1)

		if len(matches[0]) == 6 {
			stats.addStat(matches[0])
		} else {
			log.Printf("Found bad match: %+v", matches[0])
		}
	}
	return stats, err
}

func (s *realmStats) addStat(match [][]byte) {
	pktsTo, _ := strconv.Atoi(string(match[3]))
	pktsFrom, _ := strconv.Atoi(string(match[5]))
	stat := realmStat{
		Realm:     string(match[1]),
		BytesTo:   getBytes(match[2]),
		PktsTo:    pktsTo,
		BytesFrom: getBytes(match[4]),
		PktsFrom:  pktsFrom,
	}
	*s = append(*s, &stat)
}

func getBytes(s []byte) int {
	var b int
	if bytes.ContainsAny(s, "KMGT") {
		n, _ := strconv.Atoi(string(s[:len(s)-2]))
		u := s[len(s)-1]
		switch u {
		case 'K':
			b = n * (1 << 10)
		case 'M':
			b = n * (1 << 20)
		case 'G':
			b = n * (1 << 30)
		case 'T':
			b = n * (1 << 40)
		}
	} else {
		b, _ = strconv.Atoi(string(s))
	}
	return b
}
