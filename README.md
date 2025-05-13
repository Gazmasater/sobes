golangci-lint run

go install github.com/swaggo/swag/cmd/swag@latest

DROP TABLE IF EXISTS people;


internal/
└── app/
    └── mydomain/
        ├── usecase/
        │   ├── user_usecase.go        # Бизнес-логика
        │   └── user_usecase_iface.go  # Интерфейс, например UserRepository
        ├── repository/
        │   └── postgres/
        │       └── user_repository.go# Реализация интерфейса
        ├── adapters/
        │   └── http/
        │       └── handler.go         # Использует интерфейс Usecase
        └── domain.go


 curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ivan",
    "surname": "Seli",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/people/1"


curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович",
    "age": 35,
    "gender": "male",
    "nationality": "russian"
  }'


  curl -X GET http://localhost:8080/people

go test -run=NormalizeName


package yourpackage // замени на название своего пакета

import (
	"testing"
)

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"иван", "Иван"},
		{"  сЕргей", "Сергей"},
		{"ОЛЕГ  ", "Олег"},
		{"", ""},
		{"а", "А"},
		{"   ", ""},
	}

	for _, tt := range tests {
		result := NormalizeName(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeName(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}










go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.


git rm --cached textDB


curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dmitriy",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'

  curl -X DELETE "http://localhost:8080/people/5"


  curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alexey",
    "surname": "Ivanov",
    "patronymic": "Sergeevich",
    "age": 30,
    "gender": "male",
    "nationality": "ru"
  }'




package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Настройка запуска headless Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // можно false для отладки
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
			"(KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	// Переменные для хранения данных
	var title, price, rating string

	// Пример URL товара
	url := "https://www.ozon.ru/product/smartfon-samsung-galaxy-a54-128-gb-612334798/"

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),

		// Ждём появления элементов на странице
		chromedp.WaitVisible(`h1`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // на всякий случай подождем динамическую подгрузку

		// Получаем название товара
		chromedp.Text(`h1`, &title, chromedp.ByQuery),

		// Получаем цену (сначала основную, потом проверим альтернативный вариант)
		chromedp.Text(`div[data-widget="webPrice"] span`, &price, chromedp.ByQuery),

		// Получаем рейтинг, если есть
		chromedp.Text(`div[data-widget="webProductHeading"] span[style*="color:"]`, &rating, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Вывод результатов
	fmt.Println("✅ Название товара:", title)
	fmt.Println("💰 Цена:", price)
	fmt.Println("⭐ Рейтинг:", rating)
}



