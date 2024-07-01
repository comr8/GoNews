package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoNews/pkg/api"
	"GoNews/pkg/parse"
	storage "GoNews/pkg/storage"
)

func main() {
	// Чтение конфигурационного файла
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer configFile.Close()

	var config parse.Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Failed to decode config file:", err)
	}

	// Инициализация базы данных
	db := storage.InitDB()
	defer db.Close()

	// Запуск API сервера
	apiPort := "8080"
	go func() {
		err := api.StartAPI(apiPort, db)
		if err != nil {
			log.Fatal("ошибка запуска сервера апи:", err)
		}
	}()

	// Создание канала для завершения
	stopCh := make(chan struct{})

	// Обработка сигнала завершения
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		close(stopCh)
	}()

	// Запуск обхода RSS-лент
	for _, rssLink := range config.RSSLinks {
		go RunRSSParsingScheduler(rssLink, config.RequestPeriod, stopCh)
	}

	// Ожидание завершения
	<-stopCh
	fmt.Println("Работа приложения завершена")
}

// RunRSSParsingScheduler запускает планировщик для регулярного парсинга RSS-ленты
func RunRSSParsingScheduler(url string, period int, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			// Остановка планировщика
			return
		default:
			// Парсинг RSS-ленты
			posts, err := parse.ParseRSS(url)
			if err != nil {
				log.Println("оишбка парсинга RSS ленты:", err)
				continue
			}

			// Обработка полученных постов и сохранение их в бд
			for _, post := range posts {
				err := storage.NewPost(post)
				if err != nil {
					log.Println("ошибка сохранения записи в БД:", err)
				}
			}

			// Пауза перед следующим циклом парсинга
			time.Sleep(time.Duration(period) * time.Minute)
		}
	}
}
