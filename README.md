# Новостной агрегатор
### Запуск
0. У вас должна быть настроена и запущена postgresql;
1. В /pkg/storage/postgres.go в константах ```DBUser``` и ```DBPassword``` заменить значения на креды для вашей локальной БД;
2. Далее, находясь в начальной директории, где расположен go.mod (./GoNews) выполнить команду:
```go
  go run .
```
3. Открыть в браузере страницу localhost:8080 на которой будут отображены последние 10 новостей.
### Тестирование
0. ```Для запуска тестов``` из начальной директории выполнить следующую команду:
```go
go test -v ./pkg/...
```
1. ```Для получения отчета о покрытии```, из начальной директории выполнить следующую команду:

- на *Unix-like ОС:
```go
go test -v -cover -coverprofile="coverage.out" ./pkg/... && go tool cover -html="coverage.out"
```
- на Windows:
```go
go test -v -cover -coverprofile="coverage.out" ./pkg/... | go tool cover -html="coverage.out"
```
2. В директории ```/misc/brunoCollection``` расположен запрос для утилиты Bruno для отправки GET запроса с заданным количеством новостной выдачи и встроенным тестом для сверки количества получаемых новостей с количеством запрошенных (пример в скриншоте ниже).
### Результаты
0. Приложение:
![app](/misc/screenshots/app.png)
1. Покрытие тестами: 
![coverage-api](/misc/screenshots/coverage-api.png)
![coverage-storage](/misc/screenshots/coverage-postgres.png)
![coverage-parse](/misc/screenshots/coverage-parse.png)
2. Тест в Bruno:
![bruno](/misc/screenshots/count.png)