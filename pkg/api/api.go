package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	storage "GoNews/pkg/storage"

	"github.com/gorilla/mux"
)

// Структура API
type API struct {
	r  *mux.Router // Роутер для маршрутов API
	db *sql.DB     // База данных
}

// NewAPI создает новый экземпляр API с подключением к бд и настройкой роутера
func NewAPI(db *sql.DB) *API {
	api := &API{
		r:  mux.NewRouter(), // Инициализация нового роутера
		db: db,              // Инициализация подключения к бд
	}

	api.endpoints() // Настройка маршрутов API
	return api
}

// ServeHTTP позволяет API удовлетворять интерфейсу http.Handler, делегируя запросы роутеру
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.r.ServeHTTP(w, r)
}

// GetRouter возвращает роутер API для возможности его использования вне пакета
func (api *API) GetRouter() *mux.Router {
	return api.r
}

// posts обрабатывает запрос на получение последних новостей
func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		// Возвращение ошибки, если параметр 'n' не является числом
		http.Error(w, "Неверное количество новостей", http.StatusBadRequest)
		return
	}
	// Получение списка последних новостей из бд
	posts, err := storage.Posts(n)
	if err != nil {
		http.Error(w, "Не удалось получить новости", http.StatusInternalServerError)
		return
	}
	// Установка заголовка Content-Type и отправка данных в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// webAppHandler обрабатывает запросы для приложения, обслуживая статические файлы
func (api *API) webAppHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./webapp")).ServeHTTP(w, r)
}

// endpoints устанавливает маршруты API
func (api *API) endpoints() {
	// Маршрут для получения n последних новостей
	api.r.HandleFunc("/news/{n:[0-9]+}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// Маршрут для обслуживания веб-приложения
	api.r.PathPrefix("/").HandlerFunc(api.webAppHandler).Methods(http.MethodGet)
}

// StartAPI запускает API на указанном порту
func StartAPI(port string, db *sql.DB) error {
	api := NewAPI(db)
	return http.ListenAndServe(":"+port, api)
}
