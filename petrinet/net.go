package petrinet

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
)

type Net struct {
	id          string
	places      []PlaceI
	transitions []TransitionI
	dots        []string // sequence of graphviz in dot format
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
func buildDot(n *Net, t0 TransitionI) string {
	places := ""
	// Places
	for _, p := range n.places {
		toks := "\n  "
		if p.Tokens() > 0 {
			toks = "\n●" + fmt.Sprintf("%d", p.Tokens())
		}
		color := ""
		if t0 != nil && t0.isConnectedToPlace(p) {
			color = ", style=filled, fillcolor=orange"
		}
		places += "P_" + p.Id() + " [label=\"" + p.Id() + toks + "\"" + color + "]\n"
	}
	transitions := ""
	relationships := ""
	for _, t := range n.transitions {
		tp := t.(*Transition)
		// Transitions
		color := ""
		if t == t0 {
			color = ", style=filled, fillcolor=lightblue"
		}
		transitions += "T_" + t.Id() + " [label=\"" + t.Id() + "\"" + color + "]\n"
		// Relationships
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

	/* Image legend */
	graph [labeljust="l" label="%LEGEND%"]{}

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
	dot := buildDot(n, nil)
	//logger.Println(dot)

	img := dot2image(dot, map[string]string{"%LEGEND%": ""})
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		logger.Fatal(err)
	}

	fo, err := os.Create(filename)
	if err != nil {
		logger.Fatal(err)
	}
	fo.Write(buf.Bytes())
	fo.Close()
}

func (n *Net) AddAnimationFrame(t TransitionI) {
	dot := buildDot(n, t)
	//logger.Println(dot)
	n.dots = append(n.dots, dot)
}
func dot2image(dot string, params map[string]string) image.Image {
	// replace placeholders with actual values
	for k, v := range params {
		dot = strings.Replace(dot, k, v, -1)
	}

	graph, err := graphviz.ParseBytes([]byte(dot))
	if err != nil {
		logger.Fatal(err)
	}
	g := graphviz.New()
	img, err := g.RenderImage(graph)
	if err != nil {
		logger.Fatal(err)
	}
	return img
}

func (n *Net) SaveAnimationAsGif(filename string) {
	frames := make([]image.Image, len(n.dots))
	for i, dot := range n.dots {
		frames[i] = dot2image(dot, map[string]string{"%LEGEND%": fmt.Sprintf("\nFrame %d/%d", i+1, len(n.dots))})
	}

	outGif := &gif.GIF{}
	outGif.Config = image.Config{}
	for _, img := range frames {
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
