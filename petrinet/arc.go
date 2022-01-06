package petrinet

import "fmt"

/*
	[P]lace <-> [T]ransition arc (both ways)
*/
type Arc struct {
	Id     string
	Weight int
	P      PlaceI
	T      TransitionI
}

func (a *Arc) String() string {
	return fmt.Sprintf("ID [%s] Weight [%d]", a.Id, a.Weight)
}

// Used by Place to notify Transition its readiness
func (a *Arc) Notify() {
	a.T.notifyReadiness()
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
	[P]lace -> [T]ransition enable arc
*/
const undef = -1

type EnableArc struct {
	Id string
	P  PlaceI
	T  *Transition
	// weight range
	low  int // >=low
	high int // <=high
}

func NewEnableArc(id string) *EnableArc {
	arc := EnableArc{}
	arc.Id = id
	arc.low = undef
	arc.high = undef
	return &arc
}
func (a *EnableArc) IsEnabled() bool {
	toks := a.P.Tokens()
	if a.low != undef && toks < a.low {
		return false
	}
	if a.high != undef && toks > a.high {
		return false
	}
	return true
}
