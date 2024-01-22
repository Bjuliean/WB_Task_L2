package pattern

import "fmt"

// Цепочка обязанностей/цепочка вызовов(chain of responsibility) -
// поведенческий паттерн.
// Позволяет передавать запросы последовательно по цепочке, избегая
// привязки отправителя к получателю. Дает возможность обработать
// запрос сразу нескольким объектам. Уменьшает зависимость между
// клиентом и обработчиком. Реализует принцип единственной обязанности.
// Из минусов, запрос может остаться необработанным.

type Service interface {
	Handle(*Data)
	SetNext(Service)
}

type Data struct {
	GetSource bool
	Update    bool
}

//------------------------------------

type Device struct {
	Name string
	Next Service
}

func (d *Device)Handle(data *Data) {
	if data.GetSource {
		fmt.Println("already get data")
		d.Next.Handle(data)
		return
	}

	fmt.Println("getting data from device")
	data.GetSource = true
	d.Next.Handle(data)
}

func (d *Device)SetNext(service Service) {
	d.Next = service
}

//------------------------------------

type UpdateDataService struct {
	Name string
	Next Service
}

func (u *UpdateDataService)Handle(data *Data) {
	if data.Update {
		fmt.Println("data already updated")
		u.Next.Handle(data)
		return
	}

	fmt.Println("updating data from device")
	data.Update = true
	u.Next.Handle(data)
}

func (u *UpdateDataService)SetNext(service Service) {
	u.Next = service
}

//------------------------------------

type SaveDataService struct {
	Name string
	Next Service
}

func (s *SaveDataService)Handle(data *Data) {
	fmt.Println("data saved")
}

func (s *SaveDataService)SetNext(service Service) {
	s.Next = service
}