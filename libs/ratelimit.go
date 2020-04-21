package main

import (
	"container/list"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limiter struct {
	sync.RWMutex
	dataList *list.List
	dataMap  map[string]*list.Element
}

func NewLimiter() (lm *Limiter) {
	lm = &Limiter{
		dataList: list.New(),
		dataMap:  make(map[string]*list.Element),
	}
	go func() {
		for range time.Tick(time.Minute) {
			lm.gc()
		}
	}()
	return
}

func (lm *Limiter) gc() {
	var (
		item      *limiterItem
		element   *list.Element
		now       time.Time = time.Now()
		timestamp time.Time = now.Add(-5 * time.Minute)
	)
	for element = lm.dataList.Back(); element != nil; element = element.Prev() {
		item = element.Value.(*limiterItem)
		if item.lastTimestamp.After(timestamp) {
			break
		}
		if int64(now.Sub(item.lastTimestamp))/item.interval < item.burst {
			continue
		}
		// 容量满则删除
		lm.Lock()
		lm.dataList.Remove(element)
		delete(lm.dataMap, item.resource)
		lm.Unlock()
	}
}

type limiterItem struct {
	resource      string
	lastTimestamp time.Time
	interval      int64
	burst         int64
	rateLimiter   *rate.Limiter
}

func newLimiterItem(resource string, timestamp time.Time, interval int64, burst int64) (item *limiterItem) {
	item = &limiterItem{
		resource:      resource,
		lastTimestamp: timestamp,
		interval:      interval,
		burst:         burst,
		rateLimiter:   rate.NewLimiter(rate.Every(time.Duration(interval)), int(burst)),
	}
	return
}

func (item *limiterItem) update(timestamp time.Time, interval int64, burst int64) {
	if item.interval != interval {
		item.interval = interval
		item.rateLimiter.SetLimitAt(item.lastTimestamp, rate.Every(time.Duration(interval)))
	}
	if item.burst != burst {
		item.burst = burst
		item.rateLimiter.SetBurstAt(item.lastTimestamp, int(burst))
	}
	item.lastTimestamp = timestamp
}

func (lm *Limiter) GetRateLimiter(resource string, timestamp time.Time, interval int64, burst int64) (rateLimiter *rate.Limiter) {
	var (
		item    *limiterItem
		element *list.Element
		ok      bool
	)
	lm.RLock()
	element, ok = lm.dataMap[resource]
	lm.RUnlock()
	if ok {
		item = element.Value.(*limiterItem)
		item.update(timestamp, interval, burst)
		lm.dataList.MoveToFront(element)
	} else {
		item = newLimiterItem(resource, timestamp, interval, burst)
		lm.Lock()
		element = lm.dataList.PushFront(item)
		lm.dataMap[resource] = element
		lm.Unlock()
	}
	return item.rateLimiter
}
