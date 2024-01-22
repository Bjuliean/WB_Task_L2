package pattern

import "fmt"

// Стратегия(strategy) - поведенческий паттерн.
// Определяет схожие алгоритмы, и помещает каждый в
// свою отдельную структуру. Стоит использовать, когда
// нужно применить разные алгоритмы внутри одного объекта.
// Из недостатков это усложнение программы, и клиент должен
// знать в чем состоит различие стратегий, чтобы выбрать
// подходящую.

type Strategy interface {
	Route(start, end int)
}

type Navigator struct {
	Strategy
}

func (n *Navigator)SetStrategy(str Strategy) {
	n.Strategy = str
}

//------------------------------------------------

type RoadStrategy struct {
}

func (r *RoadStrategy)Route(start, end int) {
	fmt.Println(start, end, "ROAD")
}

//------------------------------------------------

type PublicTransportStrategy struct{
}

func (p *PublicTransportStrategy)Route(start, end int) {
	fmt.Println(start, end, "PUBLIC")
}

//------------------------------------------------

type WalkStrategy struct {
}

func (w *WalkStrategy)Route(start, end int) {
	fmt.Println(start, end, "WALK")
}