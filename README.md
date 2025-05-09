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
    "surname": "Selivanov",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/person/26"





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
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'


fmt.Printf("DeletePerson ID=%d\n", id)

	_, err = h.PersonRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
		return
	}

 gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ go run .
2025/05/09 13:25:56 server started on :8080
DeletePerson
2025/05/09 13:26:00 "DELETE http://localhost:8080/people/26 HTTP/1.1" from [::1]:50834 - 000 0B in 1.034772ms
2025/05/09 13:26:00 http: panic serving [::1]:50834: runtime error: slice bounds out of range [-1:]
goroutine 34 [running]:
net/http.(*conn).serve.func1()
        /usr/local/go/src/net/http/server.go:1947 +0xbe
panic({0xd6e9a0?, 0xc0001b20a8?})
        /usr/local/go/src/runtime/panic.go:787 +0x132
github.com/go-chi/chi/middleware.prettyStack.decorateFuncCallLine({}, {0xc000bca270, 0x1e}, 0x1, 0x8)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:130 +0x525
github.com/go-chi/chi/middleware.prettyStack.decorateLine({}, {0xc000bca270?, 0xced?}, 0x1, 0x8)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:106 +0x154
github.com/go-chi/chi/middleware.prettyStack.parse({}, {0xc000bc8000, 0xced, 0xc000019318?}, {0xcb40c0, 0x1cf6890})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:89 +0x3b9
github.com/go-chi/chi/middleware.PrintPrettyStack({0xcb40c0, 0x1cf6890})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:46 +0x3b
github.com/go-chi/chi/middleware.(*defaultLogEntry).Panic(0x47ff72?, {0xcb40c0?, 0x1cf6890?}, {0xc0000193e8?, 0x10000c000019478?, 0x441e00?})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/logger.go:165 +0x25
github.com/go-chi/chi/middleware.Recoverer.func1.1()
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:28 +0xc8
panic({0xcb40c0?, 0x1cf6890?})
        /usr/local/go/src/runtime/panic.go:787 +0x132
people/internal/app/people/adapters/adapterhttp.(*Handler).DeletePerson(0xc0001b64c8, {0x7c28204190f8, 0xc000220e00}, 0x2?)
        /home/gaz358/myprog/sobes/internal/app/people/adapters/adapterhttp/handlers.go:71 +0x306
net/http.HandlerFunc.ServeHTTP(0xcdb520?, {0x7c28204190f8?, 0xc000220e00?}, 0xc000013e7d?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).routeHTTP(0xc0002005a0, {0x7c28204190f8, 0xc000220e00}, 0xc00012a780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:436 +0x1f2
net/http.HandlerFunc.ServeHTTP(0xc000019870?, {0x7c28204190f8?, 0xc000220e00?}, 0xc000019850?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).ServeHTTP(0xc0002005a0, {0x7c28204190f8, 0xc000220e00}, 0xc00012a780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:70 +0x331
github.com/go-chi/chi.(*Mux).Mount.func1({0x7c28204190f8, 0xc000220e00}, 0xc00012a780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:311 +0x1c2
net/http.HandlerFunc.ServeHTTP(0xcdb520?, {0x7c28204190f8?, 0xc000220e00?}, 0xc000212507?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).routeHTTP(0xc000200540, {0x7c28204190f8, 0xc000220e00}, 0xc00012a780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:436 +0x1f2
net/http.HandlerFunc.ServeHTTP(0x7c28204190f8?, {0x7c28204190f8?, 0xc000220e00?}, 0x1cf6f01?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi/middleware.Recoverer.func1({0x7c28204190f8?, 0xc000220e00?}, 0xc000bc0300?)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:37 +0x6c
net/http.HandlerFunc.ServeHTTP(0x1cfb480?, {0x7c28204190f8?, 0xc000220e00?}, 0xc000092aa0?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi/middleware.init.0.RequestLogger.func1.1({0xefd3d0, 0xc000bc4000}, 0xc00012a640)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/logger.go:57 +0x148
net/http.HandlerFunc.ServeHTTP(0xeff2f0?, {0xefd3d0?, 0xc000bc4000?}, 0x1cf6f60?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).ServeHTTP(0xc000200540, {0xefd3d0, 0xc000bc4000}, 0xc00012a500)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:87 +0x2ef
net/http.serverHandler.ServeHTTP({0xc000bc0180?}, {0xefd3d0?, 0xc000bc4000?}, 0x1?)
        /usr/local/go/src/net/http/server.go:3301 +0x8e
net/http.(*conn).serve(0xc000b83050, {0xeff2b8, 0xc000bc0090})
        /usr/local/go/src/net/http/server.go:2102 +0x625
created by net/http.(*Server).Serve in goroutine 1
        /usr/local/go/src/net/http/server.go:3454 +0x485

