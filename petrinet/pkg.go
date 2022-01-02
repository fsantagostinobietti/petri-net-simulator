package petrinet

import (
	"context"
	"io/ioutil"
	"log"
	"os"
)

// Package level definitions

// Sync logging. WARN: use only for debugging
var logger = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds) // to disable log use  'ioutil.Discard' instead of os.Stderr
func disableLogger() {
	logger = log.New(ioutil.Discard, "", 0) // disable logs
}

// used bu sync.semaphore
var ctx = context.TODO()

// constants
var NoError error = nil
