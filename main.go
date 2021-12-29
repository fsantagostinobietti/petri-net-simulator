package main

import (
	"fmt"
	"petri-net-simulator/petrinet"
)

func main() {
	net1 := petrinet.NewNet("MyNet")
	// P1 -> T1  and  P2 -> T1
	p1 := net1.NewPlace("P1")
	p2 := net1.NewPlace("P2")
	t1 := net1.NewTransition("T1")
	const w1 = 1
	p1.ConnectTo(t1, w1)
	const w2 = 1
	p2.ConnectTo(t1, w2)
	// T1 -> Pa
	pa := net1.NewPlace("Pa")
	const wa = 1
	t1.ConnectTo(pa, wa)
	pa.SetAlertFunc(func(pi petrinet.PlaceI) bool {
		return pi.Tokens() >= 2*wa
	})

	// put tokens into net
	p1.AddTokens(3 * w1)
	p2.AddTokens(2 * w2)
	net1.AddAnimationFrame()

	// start simulation
	net1.Start()
	pa.WaitForAlert()

	net1.AddAnimationFrame()
	net1.Stop()

	// print status
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(t1)
	fmt.Println(pa)

	//net1.SavePng("mynet.png")
	net1.SaveAnimationAsGif("mynet.gif")
}
