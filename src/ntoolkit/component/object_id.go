package component

import (
	"time"
	"math/rand"
	"sync"
	"fmt"
)

var once sync.Once
var randomSeed *rand.Rand // Random number generator

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func makeObjectId() string {
	once.Do(initRand)
	return fmt.Sprintf("%d-%d", makeTimestamp(), randomSeed.Int63())
}

func initRand() {
	randomSeed = rand.New(rand.NewSource(makeTimestamp()))
}
