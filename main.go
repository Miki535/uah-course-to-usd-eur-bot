package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatId := tu.ID(update.Message.Chat.ID)
		go parse("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5", bot, chatId, "", update)
	}, th.AnyMessageWithText())

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatId := tu.ID(update.Message.Chat.ID)
		go parse("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5", bot, chatId, "true", update)
	}, th.CommandEqual("course"))

	bh.Start()
}

func parse(url string, bot *telego.Bot, chatId telego.ChatID, test string, update telego.Update) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var currencies []Currency
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		go SendMessage(bot, chatId, "Error while unmarshalling currencies!")
		log.Println(err)
	}
	for _, currency := range currencies {
		if test != "" {
			result := fmt.Sprintf("%s: Курс купівлі %s, Курс продажу %s", currency.Ccy, currency.Buy, currency.Sale)
			go SendMessage(bot, chatId, result)
		} else {
			updateText := update.Message.Text
			values, err := strconv.ParseFloat(updateText, 64)
			fff, err := strconv.ParseFloat(currency.Sale, 64)
			if err != nil {
				log.Println(err)
			}
			result := values / fff
			SendMessage(bot, chatId, fmt.Sprint(result))
		}
	}
}

func SendMessage(bot *telego.Bot, chatId telego.ChatID, text string) {
	_, err := bot.SendMessage(tu.Message(chatId, text))
	if err != nil {
		bot.SendMessage(tu.Message(chatId, "Error While Sending Message!"))
		log.Println(err)
	}
}
