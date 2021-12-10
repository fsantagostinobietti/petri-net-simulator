package petrinet

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggering(test *testing.T) {
	// build net
	p1 := NewPlace("P1")
	p2 := NewPlace("P2")
	t := NewTransition("T")
	pEnd := NewAlertPlace("PEnd")
	pEnd.AlertTokensGTE(1)
	p1.ConnectTo(t, 1)
	p2.ConnectTo(t, 2)
	t.ConnectTo(pEnd, 1)

	p1.AddTokens(1)
	p2.AddTokens(2)
	pEnd.WaitForAlert()

	assert.Equal(test, 1, pEnd.tokens())
	assert.Equal(test, 0, p1.tokens())
	assert.Equal(test, 0, p2.tokens())
}

func TestConcurrentTriggering(test *testing.T) {
	logger = log.New(ioutil.Discard, "", log.Lshortfile|log.Lmicroseconds) // disable logs
	const TRANS = 50
	const N = 1000 * TRANS
	// build petri-net
	p0 := NewPlace("P0")
	p := NewPlace("P")
	tt := make([]*Transition, TRANS)
	pEnd := NewAlertPlace("PEnd")
	pEnd.AlertTokensGTE(2 * N)
	for i := 0; i < TRANS; i++ {
		tt[i] = NewTransition("T" + fmt.Sprintf("%d", i))
		p.ConnectTo(tt[i], 1)
		p0.ConnectTo(tt[i], 1)
		tt[i].ConnectTo(pEnd, 1)
	}

	// run petri-net
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
