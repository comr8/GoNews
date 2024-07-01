package api_test

import (
	a "GoNews/pkg/api"
	storage "GoNews/pkg/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Posts(t *testing.T) {
	// Инициализация реальной бд
	storage.InitDB()
	defer storage.DB.Close()

	// Создание API
	api := a.NewAPI(storage.DB)

	// Добавление 5-ти тестовых новостей в бд
	for i := 1; i <= 5; i++ {
		testPost := storage.Post{
			Title:   fmt.Sprintf("Test Title %d", i),
			Content: "Test Content",
			PubTime: 1692644239,
			Link:    "http://example.com/test",
		}
		err := storage.NewPost(testPost)
		assert.NoError(t, err, "Failed to save post to DB")
	}

	// HTTP запрос к обработчику /news/{n}
	req, err := http.NewRequest("GET", "/news/5", nil)
	assert.NoError(t, err, "Unexpected error")

	// HTTP ResponseRecorder - регистратор ответов
	rr := httptest.NewRecorder()

	// Обработка запроса
	api.ServeHTTP(rr, req)

	// Проверка кода состояния HTTP
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	// Парсинг ответа JSON
	var responsePosts []storage.Post
	err = json.Unmarshal(rr.Body.Bytes(), &responsePosts)
	assert.NoError(t, err, "Failed to parse response JSON")

	// Последние 5-ть новостей из бд
	expectedPosts, err := storage.Posts(5)
	assert.NoError(t, err, "Failed to get latest posts from DB")

	// Проверка, что количество постов и их содержимое совпадают
	assert.Equal(t, len(expectedPosts), len(responsePosts), "Number of posts doesn't match")

	for i := 0; i < len(expectedPosts); i++ {
		assert.Equal(t, expectedPosts[i].Title, responsePosts[i].Title, "Title doesn't match")
		assert.Equal(t, expectedPosts[i].Content, responsePosts[i].Content, "Content doesn't match")
		assert.Equal(t, expectedPosts[i].PubTime, responsePosts[i].PubTime, "PubTime doesn't match")
		assert.Equal(t, expectedPosts[i].Link, responsePosts[i].Link, "Link doesn't match")
	}
}

func TestStartAPI(t *testing.T) {
	// Инициализация базы данных и мока
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка: '%s' при подключении к бд-моку", err)
	}
	defer db.Close()

	// Инициализация API с моком бд
	api := a.NewAPI(db)
	// Создание нового сервера для тестирования
	server := httptest.NewServer(api)
	defer server.Close()

	// Асинхронный вызов функции StartAPI на порту 8080
	go func() {
		if err := a.StartAPI("8080", db); err != nil {
			t.Errorf("Ошибка вызова StartAPI: %v", err)
		}
	}()

	// проверка что mock был использован так, как ожидалось
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
