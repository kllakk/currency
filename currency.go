package main

import (
	"fmt"
	"flag"
	"time"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	"io/ioutil"
)

var in = flag.String("in", "1 RUB", "Исходная валюта: Значение и RUB EUR CZK USD GEL. Например, currency -in=\"10 EUR\" ")

var client = &http.Client{Timeout: 10 * time.Second}

var types = []string {
	"RUB",
	"USD",
	"EUR",
	"CZK",
	"GEL",
}

const (
	file_name = "currency_%s.data"
	yahoo_url = "https://query.yahooapis.com/v1/public/yql?q=select+*+from+yahoo.finance.xchange+where+pair+=+%%22%s%%22&format=json&env=store%%3A%%2F%%2Fdatatables.org%%2Falltableswithkeys&callback="
)

func main() {
	flag.Parse()

	fmt.Printf("Конвертер валют %s\n", *in)

	value := strings.Split(*in, " ")

	if len(value) == 2 {
		currencyValueString := value[0]
		currencyIdString := value[1]
		if currencyValue, err := strconv.ParseFloat(currencyValueString, 64); err == nil {
			if contains(types,currencyIdString) {
				getData(currencyValue, currencyIdString)
			}
		}
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getData(v float64, id string) {
	var currencies []string
	for _, e := range types {
		s := fmt.Sprintf("%s%s", e, id)
		currencies = append(currencies, s)
	}

	url := fmt.Sprintf(yahoo_url, strings.Join(currencies,","))

	currency := new(Currency)
	file := fmt.Sprintf(file_name, id)
	getJson(url, currency, file)

	for _, e := range currency.Query.Results.Rate {
		if rate, err := strconv.ParseFloat(e.Rate, 64); err == nil {
			fmt.Printf("%s %s %f\n", e.Name, e.Rate, v / rate)
		}
	}
}

func getJson(url string, target interface{}, file string) {
	r, err := client.Get(url)
	if err != nil {
		fmt.Println("offline")
		// получить данные из файла, если он есть
		data, e := ioutil.ReadFile(file)
		if e != nil {
			fmt.Printf("Файл %s не существует\n", file)
		}

		json.Unmarshal(data, target)
	} else {
		json.NewDecoder(r.Body).Decode(target)

		data, _ := json.Marshal(target)
		ioutil.WriteFile(file, data, 0644)
	}
	defer r.Body.Close()
}

type Currency struct {
	Query struct {
		Count   int       `json:"count"`
		Created time.Time `json:"created"`
		Lang    string    `json:"lang"`
		Results struct {
			Rate []struct {
				ID   string `json:"id"`
				Name string `json:"Name"`
				Rate string `json:"Rate"`
				Date string `json:"Date"`
				Time string `json:"Time"`
				Ask  string `json:"Ask"`
				Bid  string `json:"Bid"`
			} `json:"rate"`
		} `json:"results"`
	} `json:"query"`
}
