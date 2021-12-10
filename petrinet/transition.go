package petrinet

import (
	"fmt"
)

type Transition struct {
	Id           string
	arcs_in      []*Arc
	arcs_out     []*Arc
	notification chan bool
}

// Transition constructor
func NewTransition(id string) *Transition {
	t := Transition{Id: id, notification: make(chan bool, 1)}
	t.start()
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
func (t *Transition) addOut(a *Arc) {
	t.arcs_out = append(t.arcs_out, a)
}
func lockAllInPlaces(t *Transition) {
	for {
		success := true
		max_i := 0
		for i, arc := range t.arcs_in {
			if !arc.P.trylock() {
				//logger.Printf("Transition [%s] trylock place [%s] failed!", t.Id, arc.P.Id())
				success = false
				max_i = i
				break
			}
		}
		if success {
			return // all places locked successfully
		} else {
			// unlock all places ... and try again
			for i := 0; i < max_i; i++ {
				t.arcs_in[i].P.unlock()
			}
		}
	}
}
func unlockAllInPlaces(t *Transition) {
	for _, arc := range t.arcs_in {
		arc.P.unlock()
	}
}
func consumeInTokens(t *Transition) bool {
	lockAllInPlaces(t)
	defer unlockAllInPlaces(t)

	// verify tokens can be consumed
	for _, arc := range t.arcs_in {
		if !arc.TestConsumeTokens() { // place has not enought tokens
			return false
		}
	}
	// finally consume tokens
	for _, arc := range t.arcs_in {
		arc.ConsumeTokens()
	}
	return true
}

func firingAttempt(t *Transition) bool {
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
