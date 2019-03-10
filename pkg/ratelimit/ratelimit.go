package ratelimit

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func Per(duration time.Duration, n int) rate.Limit {
	period := duration / time.Duration(n)
	return rate.Every(period)
}

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Limiter interface {
	GetVisitor(clientIP string) *rate.Limiter
	AddVisitor(clientIP string) *rate.Limiter
	CleanupVisitor(interval, expiresIn time.Duration) func(context.Context)
}

type LimiterImpl struct {
	sync.RWMutex
	sync.Once
	quit     (chan struct{})
	visitors map[string]*Visitor
	factory  func() *rate.Limiter
}

func New(limit rate.Limit, burst int) *LimiterImpl {
	return &LimiterImpl{
		quit:     make(chan struct{}),
		visitors: make(map[string]*Visitor),
		factory: func() *rate.Limiter {
			return rate.NewLimiter(limit, burst)
		},
	}
}

func (r *LimiterImpl) AddVisitor(clientIP string) *rate.Limiter {
	limiter := r.factory()
	r.Lock()
	r.visitors[clientIP] = &Visitor{limiter: limiter, lastSeen: time.Now()}
	r.Unlock()
	return limiter
}

func (r *LimiterImpl) GetVisitor(clientIP string) *rate.Limiter {
	r.RLock()
	visitor, exist := r.visitors[clientIP]
	r.RUnlock()
	if !exist {
		return r.AddVisitor(clientIP)
	}
	r.Lock()
	visitor.lastSeen = time.Now()
	r.Unlock()
	return visitor.limiter
}

func (r *LimiterImpl) CleanupVisitor(interval, expiresIn time.Duration) func(context.Context) {
	log := zap.L()
	closed := make(chan interface{})
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				r.Lock()
				for ip, v := range r.visitors {
					if time.Since(v.lastSeen) > expiresIn {
						delete(r.visitors, ip)
					}
				}
				r.Unlock()
			case <-r.quit:
				log.Info("ratelimiter closed")
				close(closed)
				return
			}
		}
	}()
	return func(ctx context.Context) {
		r.Once.Do(func() {
			close(r.quit)
		})
		select {
		case <-closed:
			log.Info("ratelimiter gracefully closed")
			return
		case <-ctx.Done():
			log.Info("ratelimiter forced closed")
			return
		}
	}
}
