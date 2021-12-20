package petrinet

import (
	"fmt"

	"github.com/golang-collections/collections/set"
)

type Transition struct {
	Id              string
	arcs_in         []*Arc
	arcs_inhibition []*InhibitionArc
	arcs_out        []*Arc
	notification    chan bool
}

// Transition constructor
func NewTransition(id string) *Transition {
	t := Transition{Id: id, notification: make(chan bool, 1)}
	return &t
}
func (t *Transition) String() string {
	s := fmt.Sprintf("Transition: ID [%s]", t.Id)
	var aa = ""
	for _, arc := range t.arcs_out {
		aa += fmt.Sprintf("%s, ", arc)
	}
	return s + " {" + aa + "}"
}
func (t *Transition) addIn(a *Arc) {
	t.arcs_in = append(t.arcs_in, a)
}
func (t *Transition) addInhibition(i *InhibitionArc) {
	t.arcs_inhibition = append(t.arcs_inhibition, i)
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
			logger.Printf("Transition [%s] lockPlaces() completed successfully!", t.Id)
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
	// verify inhibition arcs
	for _, inhibit := range t.arcs_inhibition {
		if inhibit.P.tokens() != 0 {
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
		logger.Printf("Transition [%s] ... ", t.Id)
		trigger := <-t.notification
		if !trigger {
			logger.Println("Transition stopped")
			return // stop Transition execution
		}
		if firingAttempt(t) {
			logger.Printf("Transition [%s] triggered successfully", t.Id)
		}
	}
}
func (t *Transition) start() {
	go execute(t)
	logger.Printf("Transition [%s] started ", t.Id)
}

// Stop Transition execution (blocking)
func (t *Transition) stop() {
	// send a blocking stop signal in notification channel
	t.notification <- false
}

// Used by a Place to notify to Transition it is ready for triggering (non-blocking method)
func (t *Transition) NotifyReadiness() {
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
	a.Id = fmt.Sprintf("%s >%d> %s", t.Id, weight, p.Id())
	a.Weight = weight
	a.P = p
	a.T = t
	// use arc to connect place and transition
	t.addOut(a)
	p.addIn(a)
}

func (t *Transition) InhibitedBy(p PlaceI) {
	i := new(InhibitionArc)
	i.Id = fmt.Sprintf("%s >o %s", p.Id(), t.Id)
	i.P = p
	i.T = t
	t.addInhibition(i)
}
