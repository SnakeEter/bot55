package project55

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const workerCount = 5 // количество воркеров

func main() {
	// Создаем бота с токеном
	tgBot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	// Получаем канал для обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tgBot.GetUpdatesChan(u)

	// Создаем WaitGroup для отслеживания завершения работы воркеров
	var wg sync.WaitGroup

	// Создаем фиксированное количество воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(&wg, tgBot, updates)
	}

	// Ждем завершения всех воркеров
	wg.Wait()
}

// worker — функция для обработки обновлений
func worker(wg *sync.WaitGroup, tgBot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	defer wg.Done()

	// Обрабатываем обновления
	for update := range updates {
		if update.Message != nil {
			handleUpdate(tgBot, update)
		}
	}
}

// handleUpdate — обработка сообщений от пользователей
func handleUpdate(tgBot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	messageText := update.Message.Text

	// Обработка команды /start
	if messageText == "/start" {
		err := bot.SendClassSelection(chatID, tgBot)
		if err != nil {
			log.Printf("Ошибка при отправке меню выбора класса: %v", err)
		}
	} else {
		// Обработка выбора класса и отправка расписания
		err := bot.HandleUserMessage(tgBot, update)
		if err != nil {
			log.Printf("Ошибка при обработке сообщения пользователя: %v", err)
		}
	}
}
