package pattern

import "fmt"

// Команда(command) - поведенческий паттерн.
// Инкапсулирует запрос в виде объекта, позволяя
// параметизировать клиентов для разных запросов.

type ICommand interface {
	Positive()
	Negative()
}

//-------------------------------------

type Conveyor struct {
}

func (c *Conveyor) On() {
	fmt.Println("Conveyor ON")
}

func (c *Conveyor) Off() {
	fmt.Println("Conveyor OFF")
}

func (c *Conveyor) Increase() {
	fmt.Println("Conveyor speed increased")
}

func (c *Conveyor) Decrease() {
	fmt.Println("Conveyor speed decreased")
}

//-------------------------------------

type ConveyorWorkCommand struct {
	Conveyor *Conveyor
}

func NewConveyorWorkCommand(conv *Conveyor) *ConveyorWorkCommand {
	return &ConveyorWorkCommand{
		Conveyor: conv,
	}
}

func (c *ConveyorWorkCommand) Positive() {
	c.Conveyor.On()
}

func (c *ConveyorWorkCommand) Negative() {
	c.Conveyor.Off()
}

//-------------------------------------

type ConveyorRegulatorCommand struct {
	Conveyor *Conveyor
}

func NewConveyorCommandRegulator(conv *Conveyor) *ConveyorRegulatorCommand {
	return &ConveyorRegulatorCommand{
		Conveyor: conv,
	}
}

func (c *ConveyorRegulatorCommand) Positive() {
	c.Conveyor.Increase()
}

func (c *ConveyorRegulatorCommand) Negative() {
	c.Conveyor.Decrease()
}

//-------------------------------------

type MainController struct {
	Commands map[int]*ICommand
	History  []*ICommand
}

func NewMainController() *MainController {
	return &MainController{
		Commands: make(map[int]*ICommand, 10),
		History: make([]*ICommand, 0, 10),
	}
}

func (m *MainController)SetCommand(button int, command *ICommand) {
	m.Commands[button] = command
}

func (m *MainController)PressOn(button int) {
	if v, ok := m.Commands[button]; ok {
		(*v).Positive()
		m.History = append(m.History, v)
	}
}

func (m *MainController)PressCancel(button int) {
	if v, ok := m.Commands[button]; ok {
		(*v).Negative()
		m.History = append(m.History[:button], m.History[button+1:]...)
	}
}