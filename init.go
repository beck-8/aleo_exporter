package main

import (
	"flag"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ProofRate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "proof_rate",
			Help: "The proof rate of the current machine.",
		},
	)
	blockHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "block_height",
		},
	)
	totalPuzzleRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_puzzle_rewards",
		},
	)
	estimatedNetworkSpeed15m = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "estimated_network_speed_15m",
		},
	)
	proofTarget = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "proof_target",
		},
	)
	coinbaseTarget = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "coinbase_target",
		},
	)
	networkMiners = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_miners",
		},
	)
	networkPower24h = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_power_24h",
			Help: "Accumulated power in the last 24 hours.",
		},
	)
	networkStaking = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_staking",
			Help: "Current memory usage of the service.",
		},
	)
	networkValidators = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_validators",
		},
	)
	networkDelegators = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_delegators",
		},
	)
	networkPrograms = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "network_programs",
		},
	)
	totalLogError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_log_error",
		},
	)
	totalApiError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_api_error",
		},
	)
	puzzleReward1M = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puzzle_reward_1M",
		},
	)
	puzzleRewards24h = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "puzzle_rewards_24h",
		},
	)
	zkWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "zk_workers",
		},
	)
)
var port string
var browserMetircs bool
var filePath string
var Test string

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 定义命令行参数
	flag.StringVar(&port, "port", ":8080", "Specify the port")
	flag.StringVar(&filePath, "filepath", "", "Specify aleo log file")
	flag.BoolVar(&browserMetircs, "browserMetircs", false, "Enable browserMetircs")

	// 解析命令行参数
	flag.Parse()

	if filePath == "" {
		log.Fatalln("The filepath cannot be empty")
	}

	prometheus.MustRegister(ProofRate)
	prometheus.MustRegister(totalLogError)
	if browserMetircs {
		registerMetrics(
			blockHeight,
			totalPuzzleRewards,
			estimatedNetworkSpeed15m,
			proofTarget,
			coinbaseTarget,
			networkMiners,
			networkPower24h,
			networkStaking,
			networkValidators,
			networkDelegators,
			networkPrograms,
			totalApiError,
			puzzleReward1M,
			puzzleRewards24h,
			zkWorkers,
		)
	}

}

func registerMetrics(metrics ...prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
