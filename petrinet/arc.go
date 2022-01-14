package petrinet

import "fmt"

/* Arc used to link	Place and Transition (both ways)
 */
type ArcI interface {
	Id() string
	String() string
	Place() PlaceI
	Transition() TransitionI
	// Used by Place to notify Transition its readiness
	Notify()
	// Test if it's possible consume tokens from (input) Place
	IsEnabled() bool
	// Used by Transition to remove tokens from (input) Place
	ConsumeTokens()
	// Used by Transition to add tokens to (output) Place
	FireTokens()
}

type Arc struct {
	id     string
	Weight int
	P      PlaceI
	T      TransitionI
}

func (a *Arc) Id() string {
	return a.id
}
func (a *Arc) String() string {
	return fmt.Sprintf("ID [%s] Weight [%d]", a.Id(), a.Weight)
}
func (a *Arc) Place() PlaceI {
	return a.P
}
func (a *Arc) Transition() TransitionI {
	return a.T
}
func (a *Arc) Notify() {
	if a.P.Tokens() >= a.Weight {
		a.T.notifyReadiness()
	}
}
func (a *Arc) IsEnabled() bool {
	return a.P.Tokens() >= a.Weight
}
func (a *Arc) ConsumeTokens() {
	a.P.addTokensNoLock(-a.Weight)
}
func (a *Arc) FireTokens() {
	a.P.addTokensNoLock(a.Weight)
}

/* Enable Arc type  used to link Transition to Place
Place -> Transition
*/
const undef = -1

type EnableArc struct {
	id   string
	P    PlaceI
	T    TransitionI
	low  int // weight >=low
	high int // weight <=high
}

func newEnableArc(id string) *EnableArc {
	arc := EnableArc{}
	arc.id = id
	arc.low = undef
	arc.high = undef
	return &arc
}
func (a *EnableArc) Id() string {
	return a.id
}
func (a *EnableArc) String() string {
	return fmt.Sprintf("ID [%s] Low [%d] High [%d]", a.Id(), a.low, a.high)
}
func (a *EnableArc) Place() PlaceI {
	return a.P
}
func (a *EnableArc) Transition() TransitionI {
	return a.T
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
func (a *EnableArc) ConsumeTokens() {}
func (a *EnableArc) FireTokens()    {}
func (a *EnableArc) Notify() {
	if a.IsEnabled() {
		a.T.notifyReadiness()
	}
}
