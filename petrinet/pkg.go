package petrinet

import (
	"context"
	"log"
	"os"
)

// Package level definitions

// Sync logging. WARN: use only for debugging
var logger = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds) // to disable log use  'ioutil.Discard' instead of os.Stderr

// used bu sync.semaphore
var ctx = context.TODO()
