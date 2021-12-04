package main

import (
	"fmt"
	"petri-net-simulator/petrinet"
	"time"
)

func main() {
	// P1 -> T1  and  P2 -> T1
	p1 := petrinet.NewPlace("P1")
	p2 := petrinet.NewPlace("P2")
	t1 := petrinet.NewTransition("T1")
	const w1 = 1
	p1.ConnectTo(t1, w1)
	const w2 = 1
	p2.ConnectTo(t1, w2)
	// T1 -> Pa
	pa := petrinet.NewAlertPlace("Pa")
	const wa = 1
	t1.ConnectTo(pa, wa)
	pa.AlertTokensGTE(2 * wa)

	// put tokens into net
	p1.AddTokens(3 * w1)
	p2.AddTokens(1 * w2)
	time.Sleep(100 * time.Millisecond)
	p2.AddTokens(1 * w2)

	pa.WaitForAlert()
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(t1)
	fmt.Println(pa)
}
