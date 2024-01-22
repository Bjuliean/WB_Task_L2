package pattern

import "fmt"

// Фабричный метод(factory method) - порождающий паттерн.
// Определяет общий интерфейс поведения для объектов.
// Из минусов, может привести к созданию больших параллельных
// иерархий объектов.

type Computer interface {
	PrintType()
	Details()
}

const (
	ServerType   = "server"
	PCType       = "pc"
	NotebookType = "notebook"
)

func NewComputer(typeName string) Computer {
	switch typeName {
	case ServerType:
		return NewServer()
	case PCType:
		return NewPC()
	case NotebookType:
		return NewNotebook()
	default:
		return nil
	}
}

//-------------------------------------------------

type Server struct {
	TypeName string
	Core     int
	Memory   int
}

func NewServer() Computer {
	return &Server{
		TypeName: ServerType,
		Core: 64,
		Memory: 128,
	}
}

func (s *Server)PrintType() {
	fmt.Println(s.TypeName)
}

func (s *Server)Details() {
	fmt.Println(s.Core, s.Memory)
}

//-------------------------------------------------

type PC struct {
	TypeName string
	Core     int
	Memory   int
	Monitor bool
}

func NewPC() Computer {
	return &PC{
		TypeName: PCType,
		Core: 8,
		Memory: 16,
		Monitor: true,
	}
}

func (p *PC)PrintType() {
	fmt.Println(p.TypeName)
}

func (p *PC)Details() {
	fmt.Println(p.Core, p.Memory, p.Monitor)
}

//-------------------------------------------------

type Notebook struct {
	TypeName string
	Core     int
	Memory   int
}

func NewNotebook() Computer {
	return &Notebook{
		TypeName: NotebookType,
		Core: 6,
		Memory: 16,
	}
}

func (n *Notebook)PrintType() {
	fmt.Println(n.TypeName)
}

func (n *Notebook)Details() {
	fmt.Println(n.Core, n.Memory)
}