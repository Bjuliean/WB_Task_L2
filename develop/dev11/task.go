package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

type reqInfo struct {
	UrlPath string
	Method  string
}

type Response struct {
	Status int
	Body   string
	Error  string
}

type ServerHandler struct {
	hFuncMap       map[reqInfo]func(w http.ResponseWriter, r *http.Request)
	middlewareList []func(w http.ResponseWriter, r *http.Request)
	methodList     map[string]struct{}
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		hFuncMap:       make(map[reqInfo]func(w http.ResponseWriter, r *http.Request), 10),
		middlewareList: make([]func(w http.ResponseWriter, r *http.Request), 0, 10),
		methodList: map[string]struct{}{
			http.MethodGet:    {},
			http.MethodPost:   {},
			http.MethodPut:    {},
			http.MethodDelete: {},
		},
	}
}

func (s *ServerHandler) MustAddHandleFunc(method, urlPath string, f func(w http.ResponseWriter, r *http.Request)) {
	if _, ok := s.methodList[method]; !ok {
		log.Fatal("MustAddHandleFunc: unknown method\n")
	}

	s.hFuncMap[reqInfo{
		UrlPath: urlPath,
		Method:  method,
	}] = f
}

func (s *ServerHandler) AddMiddleware(f func(w http.ResponseWriter, r *http.Request)) {
	s.middlewareList = append(s.middlewareList, f)
}

func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, v := range s.middlewareList {
		v(w, r)
	}

	if v, ok := s.hFuncMap[reqInfo{
		UrlPath: r.URL.Path,
		Method:  r.Method}]; ok {
		v(w, r)
	} else {
		enc := json.NewEncoder(w)
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   "",
			Error:  "",
		})
	}
}

func main() {
	srvHandler := NewServerHandler()

	srvHandler.AddMiddleware(Logger())

	//srvHandler.MustAddHandleFunc("GET", "/", )
	srvHandler.MustAddHandleFunc("GET", "/haha", GetHaha)

	fmt.Println("STARTING...")
	err := http.ListenAndServe("0.0.0.0:8081", srvHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func Logger() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n========LOGGER========\n")

		fmt.Printf("Method: %s\nURL: %s\nProto: %s\nAddr: %s\nUser-Agent:%s",
			r.Method, r.URL.Path, r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))

		fmt.Printf("\n======================\n")
	}
}

func GetHaha(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HAHA")
}
