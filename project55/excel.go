package project55

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func DownloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ReadExcel(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows("Sheet1") // проверяем название листа
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// кэш расписания
var scheduleCache [][]string
var lastCacheUpdate time.Time

func GetCachedSchedule(className string) ([][]string, error) {
	// обновляем кэш, если он устарел или отсутствует
	if time.Since(lastCacheUpdate) > 24*time.Hour || len(scheduleCache) == 0 {
		err := DownloadFile("https://r1.nubex.ru/s138194-7e6/f8429_f9/%D0%A0%D0%B0%D1%81%D0%BF%D0%B8%D1%81%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%83%D1%80%D0%BE%D0%BA%D0%BE%D0%B2%201%20%D1%87%D0%B5%D1%82%D0%B2%D0%B5%D1%80%D1%82%D1%8C%202024%20-%202025%20%D1%83%D1%87%D0%B5%D0%B1%D0%BD%D0%BE%D0%B3%D0%BE%20%D0%B3%D0%BE%D0%B4%D0%B0.xlsx", "data.xlsx")
		if err != nil {
			return nil, err
		}

		// читаем данные из загруженного файла
		rows, err := ReadExcel("data.xlsx")
		if err != nil {
			return nil, err
		}
		scheduleCache = rows
		lastCacheUpdate = time.Now() // обновляем время последнего обновления кэша
	}

	// фильтрация строк по классу
	var schedule [][]string
	for _, row := range scheduleCache {
		if len(row) > 0 && row[0] == className {
			schedule = append(schedule, row)
		}
	}

	// если не найдено расписание для класса
	if len(schedule) == 0 {
		return nil, nil
	}

	return schedule, nil
}

func SendSchedule(bot *tgbotapi.BotAPI, chatID int64, className string) error {
	schedule, err := GetCachedSchedule(className)
	if err != nil {
		log.Println("Ошибка при получении расписания:", err)
		return err
	}

	if len(schedule) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Расписание не найдено.")
		bot.Send(msg)
		return nil
	}

	var scheduleText string
	for _, row := range schedule {
		scheduleText += row[0] + "\n"
	}

	msg := tgbotapi.NewMessage(chatID, scheduleText)
	_, err = bot.Send(msg)
	return err
}
