package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Prometheus struct {
}

func New() *Prometheus {
	return &Prometheus{}
}

func (p *Prometheus) Run() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
