✅ 1. Установи swag (если ещё не установлен)

go install github.com/swaggo/swag/cmd/swag@latest
Убедись, что $GOPATH/bin в PATH, чтобы swag был доступен в терминале.

✅ 2. Установи необходимые пакеты

go get -u github.com/swaggo/gin-swagger         # или echo-swagger, fiber-swagger и т.д.
go get -u github.com/swaggo/files
Для gin, например:


go get -u github.com/swaggo/gin-swagger
✅ 3. Добавь комментарии для генерации документации
Пример для main.go:

go
Копировать
Редактировать
// @title           My API
// @version         1.0
// @description     This is a sample server.
// @host            localhost:8080
// @BasePath        /

func main() {
    r := gin.Default()

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run()
}
✅ 4. Сгенерируй документацию
В корне проекта:


swag init
Это создаст папку docs/ с docs.go и swagger.json.

✅ 5. Импортируй docs и подключи Swagger

import (
    "github.com/gin-gonic/gin"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "your_project/docs" // здесь твой модуль/путь к docs
)
✅ 6. Готово! Swagger доступен по адресу:

http://localhost:8080/swagger/index.html
Хочешь, я покажу пример с твоим фреймворком (Gin, Fiber, Echo и т.д.) или сгенерирую минимальный шаблон?







Вы сказали:
Убедись, что $GOPATH/bin в PATH, чтобы swag был доступен в терминале. КАК


ChatGPT сказал:
Чтобы swag был доступен в терминале, нужно добавить $GOPATH/bin в переменную окружения PATH. Вот как это сделать в зависимости от твоей системы.

🔧 Узнай значение GOPATH
Выполни в терминале:


go env GOPATH
Обычно это ~/go, если ты не задавал другое.

✅ Добавление $GOPATH/bin в PATH
Если GOPATH — это ~/go, то тебе нужно добавить ~/go/bin в PATH.

Для временного добавления (на текущую сессию):

export PATH=$PATH:$(go env GOPATH)/bin
Для постоянного добавления:
Открой файл конфигурации shell:

Если используешь Bash:


nano ~/.bashrc
Если Zsh:


nano ~/.zshrc
В конец файла добавь:


export PATH=$PATH:$(go env GOPATH)/bin
Сохрани файл и применяй изменения:


source ~/.bashrc    # или ~/.zshrc
✅ Проверь
После этого выполни:


which swag

