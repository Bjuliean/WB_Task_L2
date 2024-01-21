package pattern

import "fmt"

// Строитель(Builder) - порождающий паттерн. Разделяет
// создание сложных объектов на шаги, посредством которых
// эти объекты формируются. Дает возможность использовать
// один и тот же код строительства, для получения разных
// состояний объекта. Т.е. можно пропустить какие-то шаги,
// либо добавить еще шаг.
// Из минусов стоит отметить, что он усложняет код программы,
// из-за введения дополнительных структур, интерфейсов.
// Также объект будет привязан к конкретному строителю

type Planet struct {
	State    string
	Size     string
	Distance string
}

func (p *Planet)Info() {
	fmt.Println(p.State, p.Size, p.Distance)
}

type PlanetBuilderI interface {
	SetState(val string) PlanetBuilderI
	SetSize(val string) PlanetBuilderI
	SetDistance(val string) PlanetBuilderI
	GetPlanet() Planet
}

// -------------------------------------------

type PlanetBuilder struct {
	State    string
	Size     string
	Distance string
}

func NewPlanetBuilder() PlanetBuilderI {
	return &PlanetBuilder{}
}

func (p PlanetBuilder) SetState(val string) PlanetBuilderI {
	p.State = val
	return p
}

func (p PlanetBuilder) SetSize(val string) PlanetBuilderI {
	p.Size = val
	return p
}

func (p PlanetBuilder) SetDistance(val string) PlanetBuilderI {
	p.Distance = val
	return p
}

func (e PlanetBuilder) GetPlanet() Planet {
	return Planet{
		State: e.State,
		Size: e.Size,
		Distance: e.Distance,
	}
}