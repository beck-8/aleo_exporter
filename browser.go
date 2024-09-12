package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func browerCron(sec time.Duration) {
	ticker := time.NewTicker(sec)
	defer ticker.Stop()

	for range ticker.C {
		browerFind()
	}
}

func browerFind() {
	client := &http.Client{}

	/*
		{
		"success": true,
		"count": 0,
		"message": "",
		"Validators": 15,
		"Delegators": 51,
		"TotalBond": 1089188292907490,
		"DelegatorsBond": 1086268622170798
		}
	*/
	req1, err := http.NewRequest("GET", "https://mainnet.aleo123.io/api/v5/mainnet/validat/statistic", nil)
	if err != nil {
		log.Println("Error creating request:", err)
		totalApiError.Inc()
		return
	}

	resp1, err := client.Do(req1)
	if err != nil {
		log.Println("Error sending request:", err)
		totalApiError.Inc()
		return
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		totalApiError.Inc()
		return

	}
	var result map[string]any
	if err := json.Unmarshal(body1, &result); err != nil {
		log.Println("Error parsing JSON:", err)
		totalApiError.Inc()
		return
	}

	if d, ok := result["Validators"].(float64); ok {
		networkValidators.Set(d)
	}
	if d, ok := result["Delegators"].(float64); ok {
		networkDelegators.Set(d)
	}
	if d, ok := result["TotalBond"].(float64); ok {
		networkStaking.Set(d)
	}

	/*
		{
		"success": true,
		"count": 0,
		"message": "",
		"power": 1172661806.7255557,
		"powerd_sum": 206166463515441,
		"reward_sum": 5984151024982,
		"total_score2": "0"
		}
	*/
	req2, err := http.NewRequest("GET", "https://mainnet.aleo123.io/api/v5/mainnet/power/all", nil)
	if err != nil {
		log.Println("Error creating request:", err)
		totalApiError.Inc()
		return
	}

	resp2, err := client.Do(req2)
	if err != nil {
		log.Println("Error sending request:", err)
		totalApiError.Inc()
		return
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		totalApiError.Inc()
		return

	}

	if err := json.Unmarshal(body2, &result); err != nil {
		log.Println("Error parsing JSON:", err)
		totalApiError.Inc()
		return
	}

	if d, ok := result["power"].(float64); ok {
		estimatedNetworkSpeed15m.Set(d)
	}
	if d, ok := result["powerd_sum"].(float64); ok {
		networkPower24h.Set(d)
	}
	if d, ok := result["reward_sum"].(float64); ok {
		totalPuzzleRewards.Set(d)
	}
	if d, ok := result["reward"].(float64); ok {
		puzzleRewards24h.Set(d)
	}

	/*
		{
		  "success": true,
		  "count": 459351,
		  "message": "",
		  "block_data": [
		    {
		      "network": 1,
		      "round": 925976,
		      "height": 459350,
		      "cumulative_weight": 15686699505675434,
		      "cumulative_proof_target": 0,
		      "coinbase_target": 123792777468,
		      "proof_target": 30948194368,
		      "last_coinbase_target": 123736135561,
		      "last_coinbase_timestamp": 1724400680,
		      "timestamp": 1724400702,
		      "power": 0,
		      "reward": 0,
		      "block_reward": 23782343,
		      "TargetTotal": 0,
		      "block_hash": "ab1n4m4rlhj8unjl40flqd44gq2qhtu8w7u6cr2vyzkpeywy2uleu9smzc3tr",
		      "epoch": 1794,
		      "Transactions": 1,
		      "Solutions": 0,
		      "time": "2024-08-23T16:11:42+08:00"
		    }
		  ]
		}
	*/
	req3, err := http.NewRequest("GET", "https://mainnet.aleo123.io/api/v5/mainnet/blocks/list?page=0&page_size=1", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		totalApiError.Inc()
		return
	}

	resp3, err := client.Do(req3)
	if err != nil {
		log.Println("Error sending request:", err)
		totalApiError.Inc()
		return
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		totalApiError.Inc()
		return

	}

	if err := json.Unmarshal(body3, &result); err != nil {
		log.Println("Error parsing JSON:", err)
		totalApiError.Inc()
		return
	}

	if d, ok := result["block_data"].([]any)[0].(map[string]any)["height"].(float64); ok {
		blockHeight.Set(d)
	}
	if d, ok := result["block_data"].([]any)[0].(map[string]any)["coinbase_target"].(float64); ok {
		coinbaseTarget.Set(d)
	}
	if d, ok := result["block_data"].([]any)[0].(map[string]any)["proof_target"].(float64); ok {
		proofTarget.Set(d)
	}

	// {"success":true,"count":0,"message":"","calls":305757,"program":534,"owner":219}
	req4, err := http.NewRequest("GET", "https://mainnet.aleo123.io/api/v5/mainnet/programs/statistic", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		totalApiError.Inc()
		return
	}

	resp4, err := client.Do(req4)
	if err != nil {
		log.Println("Error sending request:", err)
		totalApiError.Inc()
		return
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp4.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		totalApiError.Inc()
		return

	}

	if err := json.Unmarshal(body4, &result); err != nil {
		log.Println("Error parsing JSON:", err)
		totalApiError.Inc()
		return
	}
	if d, ok := result["program"].(float64); ok {
		networkPrograms.Set(d)
	}

	// {"success":true,"count":0,"message":"","total":47,"new":29}
	req5, err := http.NewRequest("GET", "https://mainnet.aleo123.io/api/v5/mainnet/miner/new/day", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		totalApiError.Inc()
		return
	}

	resp5, err := client.Do(req5)
	if err != nil {
		log.Println("Error sending request:", err)
		totalApiError.Inc()
		return
	}
	defer resp5.Body.Close()

	body5, err := io.ReadAll(resp5.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		totalApiError.Inc()
		return

	}

	if err := json.Unmarshal(body5, &result); err != nil {
		log.Println("Error parsing JSON:", err)
		totalApiError.Inc()
		return
	}

	if d, ok := result["total"].(float64); ok {
		networkMiners.Set(d)
	}

	{
		req, err := http.NewRequest("GET", "https://api.aleo.info/?requestType=fetch", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			totalApiError.Inc()
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			totalApiError.Inc()
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			totalApiError.Inc()
			return

		}
		var result map[string]any
		if err := json.Unmarshal(body, &result); err != nil {
			log.Println("Error parsing JSON:", err)
			totalApiError.Inc()
			return
		}
		if d, ok := result["puzzle_reward_1M"].(float64); ok {
			puzzleReward1M.Set(d / 1000000)

		}
	}
}
