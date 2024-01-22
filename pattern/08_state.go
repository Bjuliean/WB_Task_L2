package pattern

import "fmt"

// Состояние(State) - поведенческий паттерн.
// Позволяет объекту изменять свое поведение взависимости от
// внутреннего состояния. Отлично себя покажет в системах, с
// большим количеством условных операторов

type State interface {
	NextState(*TrafficLight)
	PrevState(*TrafficLight)
}

type TrafficLight struct {
	CurState *State
}

func NewTrafficLight(st *State) TrafficLight {
	return TrafficLight{
		CurState:  st,
	}
}

func (t *TrafficLight)SetState(st *State) {
	t.CurState = st
}

func (t *TrafficLight)NextState(st *State) {
	(*st).NextState(t)
}

func (t *TrafficLight)PrevState(st *State) {
	(*st).PrevState(t)
}

//-------------------------------------------------

type GreenState struct {

}

func (g *GreenState)NextState(t *TrafficLight) {
	fmt.Println("GREEN")
	var ye State = &YellowState{}
	t.SetState(&ye)
}

func (g *GreenState)PrevState(t *TrafficLight) {
	fmt.Println("GREEN")
}

//-------------------------------------------------

type YellowState struct {

}

func (y *YellowState)NextState(t *TrafficLight) {
	fmt.Println("YELLOW")
	var rd State = &RedState{}
	t.SetState(&rd)
}

func (y *YellowState)PrevState(t *TrafficLight) {
	fmt.Println("YELLOW")
	var gr State = &GreenState{}
	t.SetState(&gr)
}

//-------------------------------------------------

type RedState struct {

}

func (r *RedState)NextState(t *TrafficLight) {
	fmt.Println("RED")
}

func (r *RedState)PrevState(t *TrafficLight) {
	fmt.Println("RED")
	var ye State = &YellowState{}
	t.SetState(&ye)
}
