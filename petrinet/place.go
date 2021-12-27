package petrinet

import (
	"fmt"

	"golang.org/x/sync/semaphore"
)

type PlaceI interface {
	String() string
	Id() string
	AddTokens(toks int) bool
	// Connect Place -> Transition with a weighted Arc
	ConnectTo(t TransitionI, weight int)
	Tokens() int
	SetAlertFunc(func(PlaceI) bool)
	// Alert generated on every change in tokens number
	SetAlertOnchange()
	WaitForAlert()

	addIn(a *Arc)
	addOut(a *Arc)
	lock()
	trylock() bool
	unlock()
	addTokensNoLock(toks int) bool
}

/*
	Place
*/
type Place struct {
	id             string
	toks           int
	sem            *semaphore.Weighted
	arcs_in        []*Arc
	arcs_out       []*Arc
	alert_onchange func(PlaceI) bool
	alert          chan bool
}

func NewPlace(id string) *Place {
	return &Place{id: id, sem: semaphore.NewWeighted(1), alert: make(chan bool, 1)}
}
func (p *Place) String() string {
	s := fmt.Sprintf("Place: ID [%s] Tokens [%d]", p.Id(), p.Tokens())
	var aa = ""
	for _, arc := range p.arcs_out {
		aa += fmt.Sprintf("%s, ", arc)
	}
	return s + " {" + aa + "}"
}
func (p *Place) Id() string {
	return p.id
}
func (p *Place) Tokens() int {
	return p.toks
}
func (p *Place) SetAlertFunc(f func(PlaceI) bool) {
	p.alert_onchange = f
}

func (p *Place) SetAlertOnchange() {
	p.SetAlertFunc(func(pi PlaceI) bool {
		return true
	})
}
func (p *Place) WaitForAlert() {
	<-p.alert
}
func (p *Place) lock() {
	p.sem.Acquire(ctx, 1)
}
func (p *Place) trylock() bool {
	return p.sem.TryAcquire(1)
}
func (p *Place) unlock() {
	p.sem.Release(1)
}

// Place notify all connected Transitions if ready for triggering
func (p *Place) notifyTransitions() {
	for _, a := range p.arcs_out {
		if p.toks >= a.Weight {
			a.Notify()
		}
	}
}
func (p *Place) generate_alert() {
	// non-blocking send
	select {
	case p.alert <- true: // alert sent
	default: // alert dropped
	}
}
func (p *Place) addTokensNoLock(toks int) bool {
	old_tokens := p.toks
	new_tokens := old_tokens + toks
	if new_tokens < 0 {
		logger.Panicf("Place [%s] cannot contain negative value for tokens", p.id)
		return false
	}
	// update tokens
	p.toks = new_tokens
	if new_tokens != old_tokens { // change in tokens
		if p.alert_onchange != nil {
			if p.alert_onchange(p) {
				// TODO generate alert
				p.generate_alert()
			}
		}
	}
	p.notifyTransitions()
	return true
}
func (p *Place) AddTokens(toks int) bool {
	p.lock()
	defer p.unlock()

	return p.addTokensNoLock(toks)
}
func (p *Place) addIn(a *Arc) {
	p.arcs_in = append(p.arcs_in, a)
}
func (p *Place) addOut(a *Arc) {
	p.arcs_out = append(p.arcs_out, a)
}
func (p *Place) ConnectTo(t TransitionI, weight int) {
	a := new(Arc)
	a.Id = fmt.Sprintf("%s >%d> %s", p.Id(), weight, t.Id())
	a.Weight = weight
	a.P = p
	a.T = t

	p.addOut(a)
	t.addIn(a)
}
