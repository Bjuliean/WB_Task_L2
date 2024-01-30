package main

import (
	"fmt"
	"wbl2/WB_Task_L2/pattern"
)

func main() {
	FacadeExample()
	BuilderExample()
	VisitorExample()
	CommandExample()
	ChainOfRespExample()
	FactoryMethodExample()
	StrategyExample()
	StateExample()
}

func FacadeExample() {
	fmt.Printf("---------------------------------\nFACADE EXAMPLE\n")

	b := &pattern.BodyManufacturer{
		Workers: 5,
	}
	e := &pattern.EngineManufacturer{
		Workers: 15,
	}
	i := &pattern.InteriorManufacturer{
		Workers: 7,
	}

	var f pattern.Facade = pattern.NewCarsFactory(b, e, i)

	f.Produce()
	fmt.Printf("\n---------------------------------\n")
}

func BuilderExample() {
	fmt.Printf("---------------------------------\nBUILDER EXAMPLE\n")

	pb := pattern.NewPlanetBuilder()

	planet := pb.SetState("empty").SetSize("small").SetDistance("far").GetPlanet()

	planet.Info()

	fmt.Printf("\n---------------------------------\n")
}

func VisitorExample() {
	fmt.Printf("---------------------------------\nVISITOR EXAMPLE\n")

	places := []pattern.Place{&pattern.Cinema{}, &pattern.Circus{}, &pattern.Zoo{}}
	t := pattern.Tourist{}

	for _, v := range places {
		v.Accept(&t)
	}

	for _, v := range t.GetVisitedPlaces() {
		fmt.Println(v)
	}

	fmt.Printf("\n---------------------------------\n")
}

func CommandExample() {
	fmt.Printf("---------------------------------\nCOMMAND EXAMPLE\n")

	conveyor := &pattern.Conveyor{}

	multipult := pattern.NewMainController()
	var cmd1 pattern.ICommand = pattern.NewConveyorWorkCommand(conveyor)
	var cmd2 pattern.ICommand = pattern.NewConveyorCommandRegulator(conveyor)
	multipult.SetCommand(0, &cmd1)
	multipult.SetCommand(1, &cmd2)

	multipult.PressOn(0)
	multipult.PressOn(1)

	multipult.PressCancel(1)
	multipult.PressCancel(0)

	fmt.Printf("\n---------------------------------\n")
}

func ChainOfRespExample() {
	fmt.Printf("---------------------------------\nCHAIN EXAMPLE\n")

	device := &pattern.Device{Name: "dev-1"}
	updater := &pattern.UpdateDataService{Name: "upd-1"}
	saver := &pattern.SaveDataService{Name: "save-1"}

	device.SetNext(updater)
	updater.SetNext(saver)

	data := pattern.Data{}
	device.Handle(&data)

	fmt.Printf("\n---------------------------------\n")
}

func FactoryMethodExample() {
	fmt.Printf("---------------------------------\nFACTORY METHOD EXAMPLE\n")

	pkg := []string{pattern.ServerType, pattern.PCType, pattern.NotebookType}

	for _, v := range pkg {
		comp := pattern.NewComputer(v)

		comp.PrintType()
		comp.Details()
	}

	fmt.Printf("\n---------------------------------\n")
}

func StrategyExample() {
	fmt.Printf("---------------------------------\nSTRATEGY EXAMPLE\n")

	pkg := []pattern.Strategy{&pattern.PublicTransportStrategy{},
		&pattern.RoadStrategy{}, &pattern.WalkStrategy{}}

	nav := &pattern.Navigator{}
	for _, v := range pkg {
		nav.SetStrategy(v)
		nav.Route(0, 10)
	}

	fmt.Printf("\n---------------------------------\n")
}

func StateExample() {
	fmt.Printf("---------------------------------\nSTATE EXAMPLE\n")

	var ye pattern.State = &pattern.YellowState{}
	tf := pattern.NewTrafficLight(&ye)

	tf.NextState(tf.CurState)
	tf.NextState(tf.CurState)
	tf.PrevState(tf.CurState)
	tf.PrevState(tf.CurState)
	tf.PrevState(tf.CurState)

	fmt.Printf("\n---------------------------------\n")
}
