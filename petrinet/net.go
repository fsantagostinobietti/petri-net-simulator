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
func (n *Net) NewPlace(id string) *Place {
	p := NewPlace(id)
	n.places = append(n.places, p)
	return p
}
func (n *Net) NewAlertPlace(id string) *AlertPlace {
	p := NewAlertPlace(id)
	n.places = append(n.places, p)
	return p
}

func (n *Net) NewTransition(id string) TransitionI {
	t := NewTransition(id)
	n.transitions = append(n.transitions, t)
	return t
}

func (n *Net) Start() {
	for _, t := range n.transitions {
		t.start()
	}
}
