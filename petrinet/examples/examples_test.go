package examples

import (
	"petri-net-simulator/petrinet"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToggleSwitch(test *testing.T) {
	const N = 4

	net := petrinet.NewNet("Test toggle switch")
	pIn, pOut := BuildToggleSwitch(net, "")
	pOut.SetAlertOnchange()

	net.Start()
	for i := 0; i < N; i++ {
		pIn.AddTokens(1)
		pOut.WaitForAlert()
		if i%2 == 0 {
			assert.Equal(test, 1, pOut.Tokens())
		} else {
			assert.Equal(test, 0, pOut.Tokens())
		}
		assert.Equal(test, 0, pIn.Tokens())
	}
	net.Stop()

	net.SavePng("00_toggle_switch.png")
}

func TestModuleNCounter(test *testing.T) {
	const N = 5

	net := petrinet.NewNet("Test module-N counter")
	pIn, pCnt := BuildModuloNCounter(net, "", N)
	net.SavePng("01_modulo_N_counter.png")

	pIn.AddTokens(2 * N)
	pIn.SetAlertFunc(func(pi petrinet.PlaceI) bool {
		return pi.Tokens() == 0
	})

	net.EnableAnimation(true)
	net.Start()
	pIn.WaitForAlert()
	assert.Equal(test, 0, pCnt.Tokens())
	net.Stop()

	net.SaveAnimationAsGif("01_modulo_N_counter.gif")
}

func TestAdder(test *testing.T) {
	net := petrinet.NewNet("Test Adder")
	pX := net.NewPlace("X")
	pY := net.NewPlace("Y")
	pRun, pSum, tNxt := BuildAdder(net, "", pX, pY)
	pNxt := net.NewPlace("Next")
	tNxt.ConnectTo(pNxt, 1)
	pNxt.SetAlertOnchange()
	net.SavePng("02_adder.png")

	const X = 3
	const Y = 2
	pX.AddTokens(X)
	pY.AddTokens(Y)

	net.EnableAnimation(true)
	net.Start()
	pRun.AddTokens(1) // run adder net
	pNxt.WaitForAlert()
	assert.Equal(test, X+Y, pSum.Tokens())
	assert.Equal(test, 0, pX.Tokens())
	assert.Equal(test, 0, pY.Tokens())
	assert.Equal(test, 0, pRun.Tokens())
	net.Stop()

	//net.SaveAnimationAsGif("02_adder.gif")
}

func TestAdder3(test *testing.T) {
	net := petrinet.NewNet("Test 3-Adder")
	pX := net.NewPlace("X")
	pY := net.NewPlace("Y")
	pZ := net.NewPlace("Z")

	pR1, pS1, tN1 := BuildAdder(net, "A1", pX, pY)
	pR2, pS2, tN2 := BuildAdder(net, "A2", pS1, pZ)
	tN1.ConnectTo(pR2, 1)

	pNxt := net.NewPlace("Next")
	tN2.ConnectTo(pNxt, 1)
	pNxt.SetAlertOnchange()
	net.SavePng("02_adder3.png")

	const X = 3
	const Y = 1
	const Z = 2
	pX.AddTokens(X)
	pY.AddTokens(Y)
	pZ.AddTokens(Z)

	net.EnableAnimation(true)
	net.Start()
	pR1.AddTokens(1) // run adder net
	pNxt.WaitForAlert()
	assert.Equal(test, X+Y+Z, pS2.Tokens())
	assert.Equal(test, 0, pX.Tokens())
	assert.Equal(test, 0, pY.Tokens())
	net.Stop()

	net.SaveAnimationAsGif("02_adder3.gif")
}
