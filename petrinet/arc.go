package petrinet

import "fmt"

/*
	[P]lace <-> [T]ransition arc (both ways)
*/
type Arc struct {
	Id     string
	Weight int
	P      PlaceI
	T      *Transition
}

func (a *Arc) String() string {
	return fmt.Sprintf("ID [%s] Weight [%d]", a.Id, a.Weight)
}

// Used by Place to notify Transition its readiness
func (a *Arc) Notify() {
	a.T.NotifyReadiness()
}

//
func (a *Arc) TestConsumeTokens() bool {
	return a.P.Tokens() >= a.Weight
}

// Used by Transition to remove tokens from (incoming) Place
func (a *Arc) ConsumeTokens() {
	a.P.addTokensNoLock(-a.Weight)
}

// Used by Transition to add tokens to (destination) Place
func (a *Arc) FireTokens() {
	a.P.addTokensNoLock(a.Weight)
}

/*
	[P]lace -> [T]ransition inhibition arc
*/
type InhibitionArc struct {
	Id string
	P  PlaceI
	T  *Transition
}
