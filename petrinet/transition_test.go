package petrinet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggering(test *testing.T) {
	// build net
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
