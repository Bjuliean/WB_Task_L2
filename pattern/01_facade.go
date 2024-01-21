package pattern

import "fmt"

// Фасад(Facade) - структурный паттерн. Его цель это
// изолировать клиентов от поведения сложной системы.
// Следовательно, стоит использовать, когда есть множество
// подсистем, у которых свои интерфейсы, и мы хотим упростить
// взаимодействие с этой группой подсистем.

// Из недостатков стоит отметить, что сам интерфейс фасада,
// в конечном итоге имеет риск стать супер-объектом, т.е. в
// программе будет слишком сильная зависимость от этого интерфейса.

type BodyProduction interface {
	ProduceBody()
}

type BodyManufacturer struct {
	Workers int
}

func (b *BodyManufacturer) ProduceBody() {
	fmt.Printf("%d workers producing body\n", b.Workers)
}

// ---------------------------------

type EngineProduction interface {
	ProduceEngine()
}

type EngineManufacturer struct {
	Workers int
}

func (e *EngineManufacturer) ProduceEngine() {
	fmt.Printf("%d workers producing engine\n", e.Workers)
}

// ---------------------------------

type InteriorProduction interface {
	ProduceInterior()
}

type InteriorManufacturer struct {
	Workers int
}

func (i *InteriorManufacturer) ProduceInterior() {
	fmt.Printf("%d workers producing interior\n", i.Workers)
}

// ---------------------------------

type Facade interface {
	Produce()
}

type CarsFactory struct {
	Body     BodyProduction
	Engine   EngineProduction
	Interior InteriorProduction
}

func NewCarsFactory(body BodyProduction, engine EngineProduction, interior InteriorProduction) *CarsFactory {
	return &CarsFactory{
		Body:     body,
		Engine:   engine,
		Interior: interior,
	}
}

func (c *CarsFactory) Produce() {
	c.Body.ProduceBody()
	c.Engine.ProduceEngine()
	c.Interior.ProduceInterior()
	fmt.Println("CAR PRODUCED")
}
