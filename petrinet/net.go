package petrinet

type Net struct {
	id          string
	places      []PlaceI
	transitions []TransitionI
}

func NewNet(id string) *Net {
	net := Net{id: id}
	return &net
}
func (n *Net) NewPlace(id string) PlaceI {
	p := newPlace(id)
	n.places = append(n.places, p)
	return p
}
func (n *Net) NewTransition(id string) TransitionI {
	t := newTransition(id)
	n.transitions = append(n.transitions, t)
	return t
}
func (n *Net) Start() {
	for _, t := range n.transitions {
		t.start()
	}
}
func (n *Net) Stop() {
	for _, t := range n.transitions {
		t.stop()
	}
}
