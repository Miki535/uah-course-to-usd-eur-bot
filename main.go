package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Currency struct {
	Ccy     string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy     string `json:"buy"`
	Sale    string `json:"sale"`
}

func main() {
	botToken := "token"

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatId := tu.ID(update.Message.Chat.ID)
		message := tu.Message(
			chatId,
			"Привіт! Щоб переглянути курс долара до гривні, або євро введіть /course",
		)
		bot.SendMessage(message)
	}, th.CommandEqual("start"))
}

func parse(url string, bot telego.Bot, chatId telego.ChatID) {
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
		//a := currency.Ccy
		b := currency.BaseCcy
		c := currency.Buy
		//d := currency.Sale
		fullMessage := fmt.Sprintf(c, b)
		bot.SendMessage(tu.Message(chatId, fullMessage))
	}
}
