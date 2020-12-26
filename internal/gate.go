package internal

type Gate struct {
	c chan struct{}
}

func NewGate(max int) *Gate {
	return &Gate{make(chan struct{}, max)}
}

func (g *Gate) Start() {
	g.c <- struct{}{}
}

func (g *Gate) Done() {
	select {
	case <-g.c:
	default:
		panic("Done called more than Start")
	}
}
