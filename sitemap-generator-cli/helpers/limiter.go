package helpers

import (
	"sync"
	"sync/atomic"
)

const (
	DEFAULT_LIMIT = int32(100)
)

type ConcurrencyLimiter struct {
	waitGroup     *sync.WaitGroup
	limit         int32
	numInProgress int32
}

// NewConcurrencyLimiter -enforce a maximum Concurrency of limit
func NewConcurrencyLimiter(limitI int, waitGroup *sync.WaitGroup) *ConcurrencyLimiter {
	limit := int32(limitI)
	if limit <= 0 {
		limit = DEFAULT_LIMIT
	}

	// allocate a limiter instance
	c := &ConcurrencyLimiter{
		limit:     limit,
		waitGroup: waitGroup,
	}
	return c
}

// Execute executes given function in go routine and increments active go routines counter
func (c *ConcurrencyLimiter) Execute(job func()) {
	c.numInProgress++
	c.waitGroup.Add(1)
	go func() {
		defer func() {
			c.waitGroup.Done()
			c.numInProgress--
		}()

		// run the job
		job()

	}()
}

// Wait - wait for all go routines to finish
// intended to be used as normal sync.WaitGroup Wait()
func (c *ConcurrencyLimiter) Wait() {
	c.waitGroup.Wait()
}

// GetNumInProgress - get a counter of how many go routines are active right now
func (c *ConcurrencyLimiter) GetNumInProgress() int32 {
	return atomic.LoadInt32(&c.numInProgress)
}
