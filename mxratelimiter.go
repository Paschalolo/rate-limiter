package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Limiter struct {
	sync.Mutex
	// bucket is filled with rate tokens per second
	rate int
	// bucket size
	bucketsize int
	//number of tokens in the bucket
	nTokens int
	// Time last token was generated
	lastToken time.Time
}

func NewLimiter(rate, limit int) *Limiter {
	return &Limiter{
		rate:       rate,
		bucketsize: limit,
		nTokens:    limit,
		lastToken:  time.Now(),
	}
}

func (L *Limiter) Wait() {
	L.Lock()
	defer L.Unlock()
	if L.nTokens > 0 {
		L.nTokens--
		return
	}
	// here there is not enough token in the bucket
	tElapsed := time.Since(L.lastToken)
	period := time.Second / time.Duration(L.rate)
	nTokens := tElapsed.Nanoseconds() / period.Nanoseconds()
	L.nTokens = int(nTokens)
	if L.nTokens > L.bucketsize {
		L.nTokens = L.bucketsize
	}
	L.lastToken = L.lastToken.Add(time.Duration(nTokens) * period)

	// we filled the cuket . There many not be enough
	if L.nTokens > 0 {
		L.nTokens--
		return
	}
	// we have to wait until more tokens are available
	// A token should be available at:
	next := L.lastToken.Add(period)
	// wait := time.Until()
	wait := time.Until(next)
	if wait >= 0 {
		time.Sleep(wait)
	}
	L.lastToken = next
}

func main() {
	limiter := NewLimiter(5, 10)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}
	time.Sleep(time.Second * 2)
	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}

}
