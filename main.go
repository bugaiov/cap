package main

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MarketCapResponse []struct {
	AvailableSupply  string `json:"available_supply"`
	HVolumeUsd       string `json:"24h_volume_usd"`
	Id               string `json:"id"`
	LastUpdated      string `json:"last_updated"`
	MarketCapUsd     string `json:"market_cap_usd"`
	Name             string `json:"name"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	PriceBtc         string `json:"price_btc"`
	PriceUsd         string `json:"price_usd"`
	Rank             string `json:"rank"`
	Symbol           string `json:"symbol"`
	TotalSupply      string `json:"total_supply"`
}

func getAndUnmarshal(url string, target interface{}) error {
	var client = &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	r, _ := client.Do(req)
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	link := "https://api.coinmarketcap.com/v1/ticker?limit=30"
	resp := new(MarketCapResponse)

	for {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Symbol",
			"Price",
			"Cap",
			"24h%",
			"7d%",
		})
		table.SetBorder(false)
		getAndUnmarshal(link, resp)
		data := [][]string{}
		print("\033[H\033[2J\n")
		for _, e := range *resp {
			marketCapUsd, _ := strconv.ParseFloat(e.MarketCapUsd, 32)
			marketCapUsd_f := marketCapUsd / 1000000
			marketCapUsd_b := strconv.FormatFloat(marketCapUsd_f, 'f', 4, 64)
			row := []string{
				e.Symbol,
				"$" + e.PriceUsd,
				marketCapUsd_b,
				e.PercentChange24h + "%",
				e.PercentChange7d + "%",
			}
			data = append(data, row)
		}
		table.AppendBulk(data)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()
		time.Sleep(60 * time.Second)
	}
}
