package petrinet

import (
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

	assert.Equal(test, 0, p1.tokens())
	assert.Equal(test, 0, p2.tokens())
}
