package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Currency struct {
	Ccy     string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy     string `json:"buy"`
	Sale    string `json:"sale"`
}

func main() {
	parse("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5")
}

func parse(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var currencies []Currency
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		log.Println(err)
	}
	for _, currency := range currencies {
		a := currency.Ccy
		b := currency.BaseCcy
		c := currency.Buy
		d := currency.Sale
		fmt.Println(a, b, c, d)
	}
}
