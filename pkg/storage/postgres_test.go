package database

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndReadFromDB(t *testing.T) {
	InitDB()

	defer DB.Close()

	// Создаем тестовый пост
	testPost := Post{
		Title:   "Test Title 2",
		Content: "Test Content",
		PubTime: 1692645688,
		Link:    "http://example.com/test",
	}

	// Сохраняем тестовый пост в базу данных
	err := NewPost(testPost)
	if err != nil {
		t.Fatalf("Не удалось сохранить новость в БД: %v", err)
	}

	// Читаем пост из базы данных по названию
	readPost, err := ReadFromDB("Test Title 2")
	if err != nil {
		t.Fatalf("Не удалось прочитать новость из БД: %v", err)
	}

	// Сравниваем значения
	if readPost.Title != testPost.Title ||
		readPost.Content != testPost.Content ||
		readPost.PubTime != testPost.PubTime ||
		readPost.Link != testPost.Link {
		t.Errorf("Сохраненные данные не соответствуют ожидаемым данным")
	}
}

func TestDeleteByTitle(t *testing.T) {
	InitDB()

	defer DB.Close()

	// Создаем тестовый пост
	testPost := Post{
		Title:   "Test Title 1",
		Content: "Test Content",
		PubTime: 1692645688,
		Link:    "http://example.com/test",
	}

	// Сохраняем тестовый пост в базу данных
	err := NewPost(testPost)
	assert.NoError(t, err, "Не удалось сохранить новость в БД")

	// Удаляем пост по названию
	err = DeleteByTitle("Test Title 3")
	assert.NoError(t, err, "Не удалось удалить новость по title")

	// Пытаемся прочитать пост с удаленным названием
	_, err = ReadFromDB("Test Title 3")
	assert.Error(t, err, "Ожидалась ошибка при попытке прочитать удаленное сообщение")
}
