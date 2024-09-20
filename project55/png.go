package project55

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUserMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	className := update.Message.Text // сообщения пользователя, где он выбрал класс

	// файл с расписанием
	schedule, err := GetScheduleForClass("file.xlsx", className)
	if err != nil {
		log.Println("Error reading schedule:", err)
		return
	}

	if len(schedule) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Расписание не найдено ВИДЕМО ШКОЛУ ВЗОРВАЛИ !!!!!!.")
		bot.Send(msg)
		return
	}

	// генерируем текст для отправки
	var scheduleText string
	for _, row := range schedule {
		scheduleText += fmt.Sprintf("%v\n", row)
	}

	// отправляем расписание пользователю
	msg := tgbotapi.NewMessage(chatID, scheduleText)
	bot.Send(msg)
}
