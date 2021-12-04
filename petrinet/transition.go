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
	// TODO
	//for {}
}
func unlockAllInPlaces(t *Transition) {
	/*for _, arc := range t.arcs_in {
		arc.P.unlock()
	}*/
}
func firingAttempt(t *Transition) bool {
	// TODO implement as an atomic operation (i.e. a transaction)
	lockAllInPlaces(t)
	for _, arc := range t.arcs_in {
		if !arc.TestConsumeTokens() { // place has not enought tokens
			return false
		}
	}
	for _, arc := range t.arcs_in {
		arc.ConsumeTokens()
	}
	unlockAllInPlaces(t)

	for _, arc := range t.arcs_out {
		arc.FireTokens()
	}
	return true
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
