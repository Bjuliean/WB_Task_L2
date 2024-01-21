package pattern

// Посетитель(Visitor) - поведенческий паттерн.
// Позволяет добавлять новые операции, не изменяя классы
// объектов, к которым применяются данные операции.
// Из недостатков стоит отметить возможное нарушение инкапсуляции

type IVisitor interface {
	Visit(Place)
}

type Place interface {
	Accept(IVisitor)
	Info() string
}

//---------------------------------------

type Zoo struct {

}

func (z *Zoo)Accept(i IVisitor) {
	i.Visit(z)
}

func (z *Zoo)Info() string {
	return "visited zoo"
}

//---------------------------------------

type Cinema struct {

}

func (c *Cinema)Accept(i IVisitor) {
	i.Visit(c)
}

func (c *Cinema)Info() string {
	return "visited cinema"
}

//---------------------------------------

type Circus struct {

}

func (c *Circus)Accept(i IVisitor) {
	i.Visit(c)
}

func (c *Circus)Info() string {
	return "visited circus"
}

//---------------------------------------

type Tourist struct {
	VisitedPlaces []string
}

func (t *Tourist)Visit(place Place) {
	t.VisitedPlaces = append(t.VisitedPlaces, place.Info())
}

func (t *Tourist)GetVisitedPlaces() []string {
	return t.VisitedPlaces
}
