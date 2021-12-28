package examples

import (
	"petri-net-simulator/petrinet"
)

/* It is used to toggle the content (zero or one token) of a output place
   using an input place as a switch.
*/
func BuildToggleSwitch(net *petrinet.Net, id string) (in petrinet.PlaceI, out petrinet.PlaceI) {
	if id != "" {
		id += "_"
	}
	pIn := net.NewPlace(id + "In")
	pOut := net.NewPlace(id + "Out")
	// On
	tOn := net.NewTransition(id + "On")
	pIn.ConnectTo(tOn, 1)
	tOn.InhibitedBy(pOut)
	tOn.ConnectTo(pOut, 1)
	// Off
	tOff := net.NewTransition(id + "Off")
	pIn.ConnectTo(tOff, 1)
	pOut.ConnectTo(tOff, 1)

	return pIn, pOut
}
