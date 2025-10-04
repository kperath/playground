package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"

	"playground/nicest-backend-framework/pokemon"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Intrumenting
// promauto automatically registers metric w/ default registry
var partyRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "pokemon_party_requests_total",
	Help: "total number of requests for party members",
},
	[]string{"party_id"},
)

var partyLastTime = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "pokemon_party_last_time_seconds",
	Help: "last time since party endpoint was called",
})

var partyLatency = promauto.NewSummary(prometheus.SummaryOpts{
	Name: "pokemon_party_latency_seconds",
	Help: "latency in seconds",
})
var partyLatencyHisto = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "pokemon_party_latency_histogram_seconds",
	Help: "latency in seconds",
})

var appInfo = promauto.NewGauge(prometheus.GaugeOpts{
	Name:        "app_info",
	ConstLabels: prometheus.Labels{"app_version": "v0.0.1"},
})

func main() {
	f, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	party, err := pokemon.NewParty(f)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.GET("/pokemon/party/:id", func(c *gin.Context) {
		start := time.Now()
		defer func() {
			dur := time.Since(start).Seconds()
			partyLatency.Observe(dur)
			partyLatencyHisto.Observe(dur)
			partyLastTime.Set(float64(time.Now().Unix()))
			id := c.Param("id")
			partyRequests.With(prometheus.Labels{"party_id": id}).Inc()
		}()
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, c.Error(err))
		}
		time.Sleep(time.Duration(rand.IntN(5)) * time.Second) // XXX: simulation code
		c.JSON(http.StatusOK, gin.H{
			"pokemon": party.Pokemon[id-1],
		})
	})
	log.Fatal(r.Run())
}
