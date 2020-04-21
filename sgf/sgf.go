package sgf

type Sgf struct {
	Index int
	Nodes map[int]*Node
}

func NewSgf() *Sgf {
	s := &Sgf{Index: 0}
	s.Nodes = make(map[int]*Node)

	return s
}
