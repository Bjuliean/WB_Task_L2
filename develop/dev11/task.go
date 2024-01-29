package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	Body   []Event
	Error  string
}

type ServerHandler struct {
	hFuncMap       map[reqInfo]func(w http.ResponseWriter, r *http.Request)
	middlewareList []func(w http.ResponseWriter, r *http.Request) error
	methodList     map[string]struct{}
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		hFuncMap:       make(map[reqInfo]func(w http.ResponseWriter, r *http.Request), 10),
		middlewareList: make([]func(w http.ResponseWriter, r *http.Request) error, 0, 10),
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

func (s *ServerHandler) AddMiddleware(f func(w http.ResponseWriter, r *http.Request) error) {
	s.middlewareList = append(s.middlewareList, f)
}

func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	for i := 0; i < len(s.middlewareList); i++ {
		err = s.middlewareList[i](w, r)
		if err != nil {
			return
		}
	}

	if v, ok := s.hFuncMap[reqInfo{
		UrlPath: r.URL.Path,
		Method:  r.Method}]; ok {
		v(w, r)
	} else {
		json.NewEncoder(w).Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  "",
		})
	}
}

const(
	PORT = ":8081"
)

func main() {
	srvHandler := NewServerHandler()

	events := NewEvents()

	srvHandler.AddMiddleware(Logger())

	srvHandler.MustAddHandleFunc("POST", "/create_event", events.CreateEvent)
	srvHandler.MustAddHandleFunc("POST", "/update_event", events.UpdateEvent)
	srvHandler.MustAddHandleFunc("POST", "/delete_event", events.DeleteEvent)
	srvHandler.MustAddHandleFunc("GET", "/events_for_day", events.EventsForDay)
	srvHandler.MustAddHandleFunc("GET", "/events_for_week", events.EventsForWeek)
	srvHandler.MustAddHandleFunc("GET", "/events_for_month", events.EventsForMonth)

	fmt.Println("STARTING...")
	log.Fatal(http.ListenAndServe("0.0.0.0"+PORT, srvHandler))
}

func Logger() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Printf("\n========LOGGER========\n")

		fmt.Printf("Method: %s\nURL: %s\nProto: %s\nAddr: %s\nUser-Agent:%s",
			r.Method, r.URL.Path, r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))

		fmt.Printf("\n======================\n")

		return nil
	}
}

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Info   string    `json:"info"`
	Date   time.Time `json:"date"`
}

type Events map[int][]Event // key = user id

func NewEvents() Events {
	return make(map[int][]Event, 10)
}

func (e *Events) CreateEvent(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	
	var ev Event
	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	for _, v := range (*e)[ev.UserID] {
		if v.ID == ev.ID {
			enc.Encode(Response{
				Status: http.StatusBadRequest,
				Body:   []Event{},
				Error:  "already exists",
			})
			return
		}
	}

	(*e)[ev.UserID] = append((*e)[ev.UserID], ev)

	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   []Event{ev},
		Error:  "",
	})
}

func (e *Events) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	
	var ev Event

	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	if _, ok := (*e)[ev.UserID]; !ok {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  "not found",
		})
		return
	}

	var tmp *Event
	for i := 0; i < len((*e)[ev.UserID]); i++ {
		if i == len((*e)[ev.UserID])-1 && (*e)[ev.UserID][i].ID != ev.ID {
			enc.Encode(Response{
				Status: http.StatusBadRequest,
				Body:   []Event{},
				Error:  "not found",
			})
			return
		}
		if (*e)[ev.UserID][i].ID == ev.ID {
			tmp = &(*e)[ev.UserID][i]
			break
		}
	}

	*tmp = ev
	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   []Event{ev},
		Error:  "",
	})
}

func (e *Events) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	
	var ev Event

	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	if _, ok := (*e)[ev.UserID]; !ok {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  "not found",
		})
		return
	}

	tmp := 0
	for i := 0; i < len((*e)[ev.UserID]); i++ {
		if i == len((*e)[ev.UserID])-1 && (*e)[ev.UserID][i].ID != ev.ID {
			enc.Encode(Response{
				Status: http.StatusBadRequest,
				Body:   []Event{},
				Error:  "not found",
			})
			return
		}
		if (*e)[ev.UserID][i].ID == ev.ID {
			tmp = i
			break
		}
	}

	(*e)[ev.UserID] = append((*e)[ev.UserID][:tmp], (*e)[ev.UserID][tmp+1:]...)
	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   []Event{},
		Error:  "",
	})
}

func (e *Events) EventsForDay(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")

	usrID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusServiceUnavailable,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	var evList []Event
	for _, v := range (*e)[usrID] {
		if v.Date.Day() == date.Day() && v.Date.Month() == date.Month() && v.Date.Year() == date.Year() {
			evList = append(evList, v)
		}
	}

	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   evList,
		Error:  "",
	})
}

func (e *Events) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")

	usrID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusServiceUnavailable,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	var evList []Event
	for _, v := range (*e)[usrID] {
		if v.Date.Day() <= date.Day() + 6 {
			evList = append(evList, v)
		}
	}

	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   evList,
		Error:  "",
	})
}

func (e *Events) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")

	usrID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusBadRequest,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		enc.Encode(Response{
			Status: http.StatusServiceUnavailable,
			Body:   []Event{},
			Error:  err.Error(),
		})
		return
	}

	var evList []Event
	for _, v := range (*e)[usrID] {
		if v.Date.Month() == date.Month() && v.Date.Year() == date.Year() {
			evList = append(evList, v)
		}
	}

	enc.Encode(Response{
		Status: http.StatusOK,
		Body:   evList,
		Error:  "",
	})
}
