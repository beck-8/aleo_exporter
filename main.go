package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go tailLogFile(filePath)
	go browerCron(30 * time.Second)

	// 暴露 Prometheus 的 metrics 接口
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Prometheus exporter is running on %s/metrics\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
