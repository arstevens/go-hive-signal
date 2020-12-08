package cache

import (
	"sync"
	"time"
)

func pollForTimedOutRecords(mutex *sync.Mutex, caches []map[string]time.Time) {
	for {
		time.Sleep(GarbageCollectionPeriod)

		mutex.Lock()
		for _, cache := range caches {
			deleteOldRecords(cache)
		}
		mutex.Unlock()
	}
}

func deleteOldRecords(cache map[string]time.Time) {
	for key, ttl := range cache {
		if ttl.Before(time.Now()) {
			delete(cache, key)
		}
	}
}
