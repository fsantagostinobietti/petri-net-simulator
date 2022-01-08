package examples

import (
	"petri-net-simulator/petrinet"
)

/* Simple adder net (i.e. SUM = X + Y), with RUN and RESET controls
 */
func BuildAdder(net *petrinet.Net, id string, pX, pY petrinet.PlaceI) (run, sum petrinet.PlaceI, nxt petrinet.TransitionI) {
	if id != "" {
		id += "_"
	}
	pSum := net.NewPlace(id + "Sum")
	pRun := net.NewPlace(id + "Run")
	// transition add X
	tAddX := net.NewTransition(id + "AddX")
	pX.ConnectTo(tAddX, 1)
	tAddX.ConnectTo(pSum, 1)
	tAddX.EnabledBy(pRun, tAddX.SetLow(1))
	// transition add Y
	tAddY := net.NewTransition(id + "AddY")
	pY.ConnectTo(tAddY, 1)
	tAddY.ConnectTo(pSum, 1)
	tAddY.EnabledBy(pRun, tAddY.SetLow(1))
	// transition next
	tNxt := net.NewTransition(id + "Next")
	pRun.ConnectTo(tNxt, 1)
	tNxt.EnabledBy(pX, tNxt.SetLow(0), tNxt.SetHigh(0))
	tNxt.EnabledBy(pY, tNxt.SetLow(0), tNxt.SetHigh(0))

	return pRun, pSum, tNxt
}
