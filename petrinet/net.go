package petrinet

import (
	"fmt"

	"github.com/goccy/go-graphviz"
)

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
func buildDot(n *Net) string {
	places := ""
	for _, p := range n.places {
		toks := ""
		if p.Tokens() > 0 {
			toks = "\n●" + fmt.Sprintf("%d", p.Tokens())
		}
		places += "P_" + p.Id() + " [label=\"" + p.Id() + toks + "\"]\n"
	}
	transitions := ""
	relationships := ""
	for _, t := range n.transitions {
		tp := t.(*Transition)
		transitions += "T_" + t.Id() + " [label=\"" + t.Id() + "\"]\n"
		for _, ain := range tp.arcs_in {
			relationships += "P_" + ain.P.Id() + " -> " + "T_" + ain.T.Id() + "\n"
		}
		for _, aen := range tp.arcs_enable {
			label := ""
			if aen.low == aen.high {
				label = "<" + fmt.Sprintf("%d", aen.low) + ">"
			} else {
				label = "<" + fmt.Sprintf("%d", aen.low) + "," + fmt.Sprintf("%d", aen.high) + ">"
			}
			relationships += "P_" + aen.P.Id() + " -> " + "T_" + aen.T.Id() + " [arrowhead=dot, label=\"" + label + "\"]\n"
		}
		for _, aout := range tp.arcs_out {
			relationships += "T_" + aout.T.Id() + " -> " + "P_" + aout.P.Id() + "\n"
		}
	}

	return `
digraph PetriNet { 

	/* Place Entities */
	{ node [shape=circle]
` + places + `
	}
	/* Transition Entities */
	{ node [shape=square]
` + transitions + `
	}
	
	/* Relationships */
` + relationships + `
}`
}

// Save Petri Net as PNG
func (n *Net) SavePng(filename string) {
	dot := buildDot(n)
	//logger.Println(dot)

	graph, err := graphviz.ParseBytes([]byte(dot))
	if err != nil {
		logger.Fatal(err)
	}
	g := graphviz.New()
	err = g.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		logger.Fatal(err)
	}
}
