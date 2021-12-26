package petrinet

import (
	"fmt"

	"golang.org/x/sync/semaphore"
)

type PlaceI interface {
	// public methods
	String() string
	Id() string
	AddTokens(toks int) bool
	// Connect Place -> Transition with a weighted Arc
	ConnectTo(t TransitionI, weight int)
	Tokens() int
	// private methods
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
	id       string
	toks     int
	sem      *semaphore.Weighted
	arcs_in  []*Arc
	arcs_out []*Arc
}

func NewPlace(id string) *Place {
	return &Place{id: id, sem: semaphore.NewWeighted(1)}
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
func (p *Place) addTokensNoLock(toks int) bool {
	new_tokens := p.toks + toks
	if new_tokens < 0 {
		logger.Panicf("Place [%s] cannot contain negative value for tokens", p.id)
		return false
	}
	p.toks = new_tokens
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

/*
	AlertPlace - a final place (with no fan-out) used to know when number of tokens is as specified
*/

type AlertPlace struct {
	id         string
	arc_in     *Arc
	toks       int
	sem        *semaphore.Weighted
	toks_alert int
	alert      chan bool
}

func NewAlertPlace(id string) *AlertPlace {
	return &AlertPlace{id: id, alert: make(chan bool, 1), sem: semaphore.NewWeighted(1)}
}
func (p *AlertPlace) String() string {
	s := fmt.Sprintf("AlertPlace: ID [%s] Tokens [%d]", p.id, p.toks)
	return s + " {}"
}
func (p *AlertPlace) Id() string {
	return p.id
}
func (p *AlertPlace) addTokensNoLock(toks int) bool {
	if toks < 0 {
		return false
	}
	p.toks += toks
	if p.toks >= p.toks_alert {
		select {
		case p.alert <- true: // message sent
		default: // message dropped
		}
	}
	return true
}
func (p *AlertPlace) AddTokens(toks int) bool {
	p.lock()
	defer p.unlock()

	return p.addTokensNoLock(toks)
}

func (p *AlertPlace) addIn(a *Arc) {
	p.arc_in = a
}

func (p *AlertPlace) addOut(a *Arc) {
	// no output for FinalPlace
}
func (p *AlertPlace) ConnectTo(t TransitionI, weight int) {
	panic("no output arcs can be connected to AlertPlace")
}

func (p *AlertPlace) AlertTokensGTE(toks int) {
	p.toks_alert = toks
}
func (p *AlertPlace) WaitForAlert() {
	<-p.alert
}
func (p *AlertPlace) Tokens() int {
	return p.toks
}
func (p *AlertPlace) lock() {
	p.sem.Acquire(ctx, 1)
}
func (p *AlertPlace) trylock() bool {
	return p.sem.TryAcquire(1)
}
func (p *AlertPlace) unlock() {
	p.sem.Release(1)
}
