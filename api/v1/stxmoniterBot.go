package v1

import (
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

func StxmoniterBot() {
	bot, err := tgbotapi.NewBotAPI("7463838551:AAGqKaxSB-_Pu8XifNJKWAftaUfPli809Mg")
	if err != nil {
		log.Panic(err)
	}
	checkBlockChanges(bot)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// 监听命令
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "您的机器人已启动！")
				bot.Send(msg)
			}
		}
	}
}

func checkBlockChanges(bot *tgbotapi.BotAPI) {
	url := "https://explorer.hiro.so/?chain=mainnet"
	previousValue := ""

	for {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Println(err)
			return
		}

		doc.Find("#recent-blocks").Each(func(i int, s *goquery.Selection) {
			recentBlocks := s.Text()

			if previousValue == "" {
				previousValue = recentBlocks
			} else if recentBlocks != previousValue {
				sendNotification(bot, recentBlocks)
				previousValue = recentBlocks
			}
		})

		time.Sleep(60 * time.Second) // 每60秒检查一次
	}
}

func sendNotification(bot *tgbotapi.BotAPI, message string) {
	chatID := int64(7463838551)

	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
