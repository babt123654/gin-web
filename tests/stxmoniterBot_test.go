package tests

import (
	v1 "gin-web/api/v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"testing"
)

func TestStxmoniterBot(t *testing.T) {
	//v1.StxmoniterBot()
	bot, err := tgbotapi.NewBotAPI("7463838551:AAGqKaxSB-_Pu8XifNJKWAftaUfPli809Mg")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	// 调用发送消息的方法
	err = v1.SendMessage(bot, 7463838551, "这是一条消息提示！")
	if err != nil {
		log.Fatal(err)
	}
}
