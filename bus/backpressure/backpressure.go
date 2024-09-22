package backpressure

import (
	"sync"
	"time"
	"fmt"
)

// BackPressureManager handles scenarios where publishers outpace subscribers
type BackPressureManager struct {
	topicRates map[string]*rate
	mutex      sync.RWMutex
	threshold int
}

type rate struct {
	count     int
	timestamp time.Time
}

// NewBackPressureManager creates a new BackPressureManager
func NewBackPressureManager(threshold int) *BackPressureManager {
	return &BackPressureManager{
		topicRates: make(map[string]*rate),
		threshold: threshold,
	}
}

// CheckPressure checks if a new message can be published to a topic
func (bpm *BackPressureManager) CheckPressure(topic string) error {
    bpm.mutex.Lock()
    defer bpm.mutex.Unlock()

    now := time.Now()
    if r, exists := bpm.topicRates[topic]; exists {
        if now.Sub(r.timestamp) < time.Second {
            r.count++
            if r.count > bpm.threshold {
                return fmt.Errorf("back pressure applied: too many messages for topic %s", topic)
            }
        } else {
            r.count = 1
            r.timestamp = now
        }
    } else {
        bpm.topicRates[topic] = &rate{count: 1, timestamp: now}
    }
    
    return nil
}
