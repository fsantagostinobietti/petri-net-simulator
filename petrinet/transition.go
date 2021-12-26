package petrinet

import (
	"fmt"

	"github.com/golang-collections/collections/set"
)

type TransitionI interface {
	Id() string
	String() string
	ConnectTo(p PlaceI, weight int)
	EnabledBy(p PlaceI, params ...func(*EnableArc))
	SetLow(low int) func(*EnableArc)
	SetHigh(high int) func(*EnableArc)
	InhibitedBy(p PlaceI)
	// private
	notifyReadiness()
	addIn(a *Arc)
	start()
}

type Transition struct {
	id           string
	arcs_in      []*Arc
	arcs_enable  []*EnableArc
	arcs_out     []*Arc
	notification chan bool
}

// Transition constructor
func NewTransition(id string) TransitionI {
	t := Transition{id: id, notification: make(chan bool, 1)}
	return &t
}
func (t *Transition) Id() string {
	return t.id
}
func (t *Transition) String() string {
	s := fmt.Sprintf("Transition: ID [%s]", t.Id())
	var aa = ""
	for _, arc := range t.arcs_out {
		aa += fmt.Sprintf("%s, ", arc)
	}
	return s + " {" + aa + "}"
}
func (t *Transition) addIn(a *Arc) {
	t.arcs_in = append(t.arcs_in, a)
}
func (t *Transition) addEnableArc(e *EnableArc) {
	t.arcs_enable = append(t.arcs_enable, e)
}
func (t *Transition) addOut(a *Arc) {
	t.arcs_out = append(t.arcs_out, a)
}

/* Locks all In and Out transition places
 */
func lockPlaces(t *Transition, places *set.Set) {

	for {
		locked := make([]PlaceI, 0, places.Len())

		// try locking places
		success := true
		places.Do(func(i interface{}) {
			if success {
				place := i.(PlaceI)
				if place.trylock() {
					locked = append(locked, place)
				} else {
					success = false
				}
			}
		})

		if success {
			logger.Printf("Transition [%s] lockPlaces() completed successfully!", t.Id())
			return // all places locked successfully
		} else {
			// unlock all places ... and try again
			for _, place := range locked {
				place.unlock()
			}
		}
	}
}

func unlockPlaces(t *Transition, places *set.Set) {
	places.Do(func(i interface{}) {
		place := i.(PlaceI)
		place.unlock()
	})
}
func consumeInTokens(t *Transition) bool {
	// verify tokens can be consumed
	for _, arc := range t.arcs_in {
		if !arc.TestConsumeTokens() { // input place has not enought tokens
			return false
		}
	}
	// verify enabling arcs
	for _, enable := range t.arcs_enable {
		if !enable.IsEnabled() {
			return false
		}
	}
	// finally consume tokens
	for _, arc := range t.arcs_in {
		arc.ConsumeTokens()
	}
	return true
}
func uniquePlaces(t *Transition) *set.Set {
	arcs := make([]*Arc, 0, len(t.arcs_in)+len(t.arcs_out))
	arcs = append(arcs, t.arcs_in...)
	arcs = append(arcs, t.arcs_out...)

	uniques := set.New()
	for _, a := range arcs {
		uniques.Insert(a.P)
	}
	return uniques
}
func firingAttempt(t *Transition) bool {
	all_places := uniquePlaces(t)
	lockPlaces(t, all_places)
	defer unlockPlaces(t, all_places)

	ok := consumeInTokens(t)
	if ok {
		for _, arc := range t.arcs_out {
			arc.FireTokens()
		}
	}
	return ok
}
func execute(t *Transition) {
	for {
		logger.Printf("Transition [%s] ... ", t.Id())
		trigger := <-t.notification
		if !trigger {
			logger.Println("Transition stopped")
			return // stop Transition execution
		}
		if firingAttempt(t) {
			logger.Printf("Transition [%s] triggered successfully", t.Id())
		}
	}
}
func (t *Transition) start() {
	go execute(t)
	logger.Printf("Transition [%s] started ", t.Id())
}

// Stop Transition execution (blocking)
func (t *Transition) stop() {
	// send a blocking stop signal in notification channel
	t.notification <- false
}

// Used by a Place to notify to Transition it is ready for triggering (non-blocking method)
func (t *Transition) notifyReadiness() {
	// async write
	select {
	case t.notification <- true:
		// message sent
	default:
		// message dropped
	}
}
func (t *Transition) ConnectTo(p PlaceI, weight int) {
	// create arc
	a := new(Arc)
	a.Id = fmt.Sprintf("%s >%d> %s", t.Id(), weight, p.Id())
	a.Weight = weight
	a.P = p
	a.T = t
	// use arc to connect place and transition
	t.addOut(a)
	p.addIn(a)
}

func (t *Transition) SetLow(low int) func(*EnableArc) {
	return func(a *EnableArc) {
		a.low = low
	}
}
func (t *Transition) SetHigh(high int) func(*EnableArc) {
	return func(a *EnableArc) {
		a.high = high
	}
}
func (t *Transition) EnabledBy(p PlaceI, params ...func(*EnableArc)) {
	e := new(EnableArc)
	e.Id = fmt.Sprintf("%s >● %s", p.Id(), t.Id())
	e.P = p
	e.T = t
	// set params
	for _, f := range params {
		f(e)
	}
	t.addEnableArc(e)
}

func (t *Transition) InhibitedBy(p PlaceI) {
	t.EnabledBy(p, t.SetLow(0), t.SetHigh(0))
}
