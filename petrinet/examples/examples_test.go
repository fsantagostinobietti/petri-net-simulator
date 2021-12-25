package examples

import (
	"petri-net-simulator/petrinet"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToggleSwitch(test *testing.T) {
	const N = 4

	net := petrinet.NewNet("Test toggle switch")
	pIn, pOut := BuildToggleSwitch(net, "toggle")

	net.Start()
	for i := 0; i < N; i++ {
		pIn.AddTokens(1)
		time.Sleep(50 * time.Millisecond)
		if i%2 == 0 {
			assert.Equal(test, 1, pOut.Tokens())
		} else {
			assert.Equal(test, 0, pOut.Tokens())
		}
		assert.Equal(test, 0, pIn.Tokens())
	}
}