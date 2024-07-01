package parse

import (
	storage "GoNews/pkg/storage"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

// Структура для конфигурации парсера
type Config struct {
	RSSLinks      []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

// CleanHTMLTags удаляет HTML-теги из текста и возвращает очищенный текст
func CleanHTMLTags(input string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	var cleanedText string

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return cleanedText
		case html.TextToken:
			text := tokenizer.Token().Data
			cleanedText += text
		}
	}
}

// ParseRSS выполняет парсинг RSS-ленты по указанному URL и возвращает список постов
func ParseRSS(url string) ([]storage.Post, error) {
	var posts []storage.Post

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	for _, item := range feed.Items {
		// Парсинг времени публикации с использованием заданного формата
		pubTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		if err != nil {
			return nil, err
		}

		post := storage.Post{
			Title:   item.Title,
			Content: CleanHTMLTags(item.Description),
			PubTime: pubTime.Unix(),
			Link:    item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// ParseRSSFixture парсит фикстурный XML RSS-потока и возвращает список новостей
func ParseRSSFixture(fixtureXML string) ([]storage.Post, error) {
	var posts []storage.Post

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(fixtureXML)
	if err != nil {
		return nil, err
	}

	for _, item := range feed.Items {
		pubTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		if err != nil {
			return nil, err
		}

		post := storage.Post{
			Title:   item.Title,
			Content: CleanHTMLTags(item.Description),
			PubTime: pubTime.Unix(),
			Link:    item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
