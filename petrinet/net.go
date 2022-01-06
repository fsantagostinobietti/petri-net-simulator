package petrinet

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
)

type Net struct {
	id           string
	places       []PlaceI
	transitions  []TransitionI
	animation    bool // enable/disable animation recording
	animationSem chan bool
	frames       []frame // animation sequence in graphviz/dot format
}

type frame struct {
	dot   string
	delay int
}

func NewNet(id string) *Net {
	net := Net{id: id, animationSem: make(chan bool, 1)}
	net.animationSem <- true
	return &net
}
func (n *Net) NewPlace(id string) PlaceI {
	p := newPlace(id)
	n.places = append(n.places, p)
	return p
}
func (n *Net) NewTransition(id string) TransitionI {
	t := newTransition(n, id)
	n.transitions = append(n.transitions, t)
	return t
}
func (n *Net) Start() {
	// initial frame
	dot := buildDot(n, nil)
	n.addAnimationFrame([]frame{{dot, 200}})

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
			low, high := aen.low, aen.high
			if low == high && low != undef {
				label = "<" + fmt.Sprintf("%d", aen.low) + ">"
			} else {
				label += "<"
				if low != undef {
					label += fmt.Sprintf("%d", low)
				}
				label += ","
				if high != undef {
					label += fmt.Sprintf("%d", high)
				}
				label += ">"
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
func (n *Net) SavePng(filename string) error {
	dot := buildDot(n, nil)
	//logger.Println(dot)

	img := dot2image(dot, map[string]string{"%LEGEND%": ""})
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return err
	}

	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()
	fo.Write(buf.Bytes())
	return NoError
}

func (n *Net) addAnimationFrame(frames []frame) {
	if n.animation {
		for _, frame := range frames {
			<-n.animationSem // wait for green light
			n.frames = append(n.frames, frame)
			n.animationSem <- true // release semaphore
		}
	}
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

func (n *Net) EnableAnimation(enable bool) {
	n.animation = enable
}

func (n *Net) SaveAnimationAsGif(filename string) error {
	if !n.animation {
		return fmt.Errorf("SaveAnimationAsGif() failed for [%s]! Enable animation first", n.id)
	}
	imgs := make([]image.Image, len(n.frames)+1)
	delays := make([]int, len(n.frames)+1)
	width, height := 0, 0
	for i, frame := range n.frames {
		img := dot2image(frame.dot, map[string]string{"%LEGEND%": fmt.Sprintf("\nFrame %d/%d", i+1, len(n.frames))})
		// adjust max width/height
		if img.Bounds().Max.X > width {
			width = img.Bounds().Max.X
		}
		if img.Bounds().Max.Y > height {
			height = img.Bounds().Max.Y
		}
		imgs[i+1] = img
		delays[i+1] = frame.delay
	}
	// start animatin with a blanck image
	blankImg := image.NewPaletted(image.Rect(0, 0, width, height), color.Palette([]color.Color{color.White}))
	imgs[0] = blankImg
	delays[0] = 200

	outGif := &gif.GIF{}
	outGif.Config = image.Config{Width: width, Height: height}
	for i, img := range imgs {
		// convert to paletted image
		palettedImage := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, img, img.Bounds().Min, draw.Over)

		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, delays[i]) // 100ths of a second
	}
	// save to file
	f, _ := os.Create(filename)
	defer f.Close()
	gif.EncodeAll(f, outGif)

	return NoError
}
