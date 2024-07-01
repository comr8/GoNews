package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Константы для настройки подключения к базе данных
const (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "postgres"
	DBPassword = "password"
	DBName     = "postgres"
)

var DB *sql.DB

// Инициализация базы данных
func InitDB() *sql.DB {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	DB, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}

// Сохранение новости в базе данных
func NewPost(post Post) error {
	_, err := DB.ExecContext(context.Background(), `
		INSERT INTO news (title, content, pubtime, link)
		VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		post.Title,
		post.Content,
		post.PubTime,
		post.Link,
	)
	if err != nil {
		return fmt.Errorf("ошибка добавления новой публикации: %v", err)
	}
	return nil
}

// Чтение новости из базы данных по названию
func ReadFromDB(title string) (Post, error) {
	var post Post

	query := `
		SELECT id, title, content, pubtime, link
		FROM news
		WHERE title = $1
	`
	row := DB.QueryRow(query, title)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return post, err
	}

	return post, nil
}

// Получение n последних новостей из базы данных
func Posts(limit int) ([]Post, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := DB.QueryContext(context.Background(), `
		SELECT
			id,
			title,
			content,
			pubTime,
			link
		FROM news
		ORDER BY pubTime ASC
		LIMIT $1;
	`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из таблицы: %v", err)
	}
	var posts []Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Post
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.PubTime,
			&t.Link,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строк: %v", err)
		}
		// добавление переменной в массив результатов
		posts = append(posts, t)

	}
	return posts, rows.Err()
}

// Удаление новости из базы данных по названию
func DeleteByTitle(title string) error {
	_, err := DB.Exec("DELETE FROM news WHERE title = $1", title)
	return err
}
