package petrinet

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func adderRoutine(wg *sync.WaitGroup, p *Place, iters int) {
	defer wg.Done()
	for ; iters > 0; iters-- {
		p.AddTokens(1)
	}
}
func TestAddTokensConcurrency(t *testing.T) {
	p := newPlace("P")
	wg := sync.WaitGroup{}

	const N = 100
	const TOKS = 10000
	wg.Add(N)
	for i := 0; i < N; i++ {
		go adderRoutine(&wg, p, TOKS)
	}

	wg.Wait()
	assert.Equal(t, N*TOKS, p.Tokens())
}
