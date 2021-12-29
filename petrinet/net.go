package petrinet

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"

	"github.com/goccy/go-graphviz"
)

type Net struct {
	id          string
	places      []PlaceI
	transitions []TransitionI
	frames      []image.Image // animation frames
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

// build net graph as graphviz dot text
func buildDot(n *Net) string {
	places := ""
	for _, p := range n.places {
		toks := "\n  "
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

func (n *Net) AddAnimationFrame() {
	dot := buildDot(n)
	//logger.Println(dot)

	graph, err := graphviz.ParseBytes([]byte(dot))
	if err != nil {
		logger.Fatal(err)
	}
	g := graphviz.New()
	img, err := g.RenderImage(graph)
	if err != nil {
		logger.Fatal(err)
	}
	n.frames = append(n.frames, img)
}

func (n *Net) SaveAnimationAsGif(filename string) {
	outGif := &gif.GIF{}
	outGif.Config = image.Config{}
	for _, img := range n.frames {
		// convert image to paletted
		palettedImage := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, img, img.Bounds().Min, draw.Over)

		// adjust max width/height
		if img.Bounds().Max.X > outGif.Config.Width {
			outGif.Config.Width = img.Bounds().Max.X
		}
		if img.Bounds().Max.Y > outGif.Config.Height {
			outGif.Config.Height = img.Bounds().Max.Y
		}

		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 100) // 100ths of a second
	}
	// save to file
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
