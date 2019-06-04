package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String(
	"listen-address",
	":8080",
	"The address to listen on for HTTP requests.",
)
var timer = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "timer",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0},
		Help: "Front-end timer metrics.",
	}, []string{"name", "api", "status"},
)

type Metrics struct {
	Type string `json:"type"`
	Labels []string `json:"labels"`
	Value float64 `json:"value"`
}

func receive(w http.ResponseWriter, r *http.Request) {
	var m Metrics
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		log.Printf("Processed error: %s", err)
		http.Error(w, fmt.Sprintf("Get error: %s", err), 400)
		return
	}

	switch m.Type {
	case "timer":
		observer := timer.WithLabelValues(m.Labels...)
		observer.Observe(m.Value)
	}
}

func main(){
	flag.Parse()
	http.HandleFunc("/api", receive)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
