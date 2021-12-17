package petrinet

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTriggering(test *testing.T) {
	/* build net:
	(P1) -1-
	        \
	         *-> [T] -1-> (PEnd)
	        /
	(P2) -2-
	*/
	net := NewNet("TestNet")
	p1 := net.NewPlace("P1")
	p2 := net.NewPlace("P2")
	t := net.NewTransition("T")
	pEnd := net.NewAlertPlace("PEnd")
	pEnd.AlertTokensGTE(1)
	p1.ConnectTo(t, 1)
	p2.ConnectTo(t, 2)
	t.ConnectTo(pEnd, 1)

	p1.AddTokens(1)
	p2.AddTokens(2)
	net.Start()
	pEnd.WaitForAlert()

	assert.Equal(test, 1, pEnd.tokens())
	assert.Equal(test, 0, p1.tokens())
	assert.Equal(test, 0, p2.tokens())
}

// test simple net with a transition
// using same place for both input and output
func TestCloseLoopTriggering(test *testing.T) {
	disableLogger()
	/* build net:
	(P0) -2-> [T1] -1-> (PEnd)
	  ^------1-'
	*/
	const N = 16
	// build petri-net
	net := NewNet("TestNet")
	p0 := net.NewPlace("P0")
	t1 := net.NewTransition("T1")
	pEnd := net.NewAlertPlace("PEnd")
	p0.ConnectTo(t1, 2)
	t1.ConnectTo(p0, 1)
	t1.ConnectTo(pEnd, 1)
	pEnd.AlertTokensGTE(powInt(2, N) - 1)

	// run net
	p0.AddTokens(powInt(2, N))
	net.Start()
	pEnd.WaitForAlert()

	assert.Equal(test, powInt(2, N)-1, pEnd.tokens())
}

// Test multiple transitions triggering concurrently againt the same places
func TestConcurrentTriggering(test *testing.T) {
	disableLogger()
	const TRANS = 50
	const N = 1000 * TRANS
	// build petri-net
	net := NewNet("TestNet")
	p0 := net.NewPlace("P0")
	p := net.NewPlace("P")
	tt := make([]*Transition, TRANS)
	pEnd := net.NewAlertPlace("PEnd")
	pEnd.AlertTokensGTE(2 * N)
	for i := 0; i < TRANS; i++ {
		tt[i] = net.NewTransition("T" + fmt.Sprintf("%d", i))
		p.ConnectTo(tt[i], 1)
		p0.ConnectTo(tt[i], 1)
		tt[i].ConnectTo(pEnd, 1)
	}

	// run petri-net
	net.Start()
	p0.AddTokens(2 * N)
	p.AddTokens(N)
	for i := 0; i < N; i++ {
		p.AddTokens(1)
	}
	pEnd.WaitForAlert()

	assert.Equal(test, 2*N, pEnd.tokens())
	assert.Equal(test, 0, p.tokens())
	assert.Equal(test, 0, p0.tokens())
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// Test transition atomicity.
// Without it a deadlock will happen during test.
func TestAtomicTriggering(test *testing.T) {
	disableLogger()

	const N = 16
	// build petri-net
	net := NewNet("Avoid Deadlock")
	p1 := net.NewPlace("P1")
	p2 := net.NewPlace("P2")

	t1 := net.NewTransition("T1")
	p1.ConnectTo(t1, 1) //
	t1.ConnectTo(p2, 1) // P1 -(1)-> T1 -(1)-> P2

	t2 := net.NewTransition("T2")
	p2.ConnectTo(t2, 2) // T2 consume 1 token every time it triggers:
	t2.ConnectTo(p1, 1) // P2 -(2)-> T2 -(1)-> P1

	palert := net.NewAlertPlace("Alert")
	palert.AlertTokensGTE(powInt(2, N) - 1)
	t2.ConnectTo(palert, 1) // T2 -(1)-> Palert

	// run petri net
	p1.AddTokens(powInt(2, N))
	net.Start()
	palert.WaitForAlert()

	assert.Equal(test, powInt(2, N)-1, palert.tokens()) //2^N-1
}
