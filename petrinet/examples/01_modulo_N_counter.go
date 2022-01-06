package examples

import (
	"petri-net-simulator/petrinet"
)

/* Modulo N counter.
   Adding a token to 'in' place increases 'cnt' value. When N-1 is reached, 'cnt' valule is reset to 0.
*/
func BuildModuloNCounter(net *petrinet.Net, id string, N int) (in petrinet.PlaceI, cnt petrinet.PlaceI) {
	if id != "" {
		id += "_"
	}
	pIn := net.NewPlace(id + "In")
	pCnt := net.NewPlace(id + "Cnt")
	// Increment
	tInc := net.NewTransition(id + "Inc")
	pIn.ConnectTo(tInc, 1)
	tInc.ConnectTo(pCnt, 1)
	tInc.EnabledBy(pCnt, tInc.SetHigh(N-2))
	// Reset
	tRst := net.NewTransition(id + "Rst")
	pCnt.ConnectTo(tRst, N-1)
	pIn.ConnectTo(tRst, 1)

	return pIn, pCnt
}
