package chapter5

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type APIConnection interface {
	ReadFile(context.Context) error
	ResolveAddress(context.Context) error
}

type NoRateLimitAPIConnection struct{}

func (a *NoRateLimitAPIConnection) ReadFile(c context.Context) error {
	return nil
}

func (a *NoRateLimitAPIConnection) ResolveAddress(c context.Context) error {
	return nil
}

func request(apiConnection APIConnection) {
	defer log.Printf("Done.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}

			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}

func NoRateLimit() {
	open := func() *NoRateLimitAPIConnection {
		return &NoRateLimitAPIConnection{}
	}

	apiConnection := open()
	request(apiConnection)
}

type SimpleRateLimitAPIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *SimpleRateLimitAPIConnection) ReadFile(c context.Context) error {
	if err := a.rateLimiter.Wait(c); err != nil {
		return err
	}

	return nil
}

func (a *SimpleRateLimitAPIConnection) ResolveAddress(c context.Context) error {
	if err := a.rateLimiter.Wait(c); err != nil {
		return err
	}

	return nil
}

func SimpleRateLimit() {
	apiConnection := &SimpleRateLimitAPIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
	request(apiConnection)
}

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func (l *multiLimiter) Wait(c context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(c); err != nil {
			return err
		}
	}

	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

type MultiRateLimitAPIConnection struct {
	rateLimiter RateLimiter
}

func (a *MultiRateLimitAPIConnection) ReadFile(c context.Context) error {
	if err := a.rateLimiter.Wait(c); err != nil {
		return err
	}
	return nil
}

func (a *MultiRateLimitAPIConnection) ResolveAddress(c context.Context) error {
	if err := a.rateLimiter.Wait(c); err != nil {
		return err
	}
	return nil
}

func MultiRateLimit() {
	open := func() *MultiRateLimitAPIConnection {
		secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
		minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

		return &MultiRateLimitAPIConnection{
			rateLimiter: MultiLimiter(secondLimit, minuteLimit),
		}
	}

	request(open())
}
