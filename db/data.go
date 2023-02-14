package db

import "sync"

var (
	cachedDB    AuthDB
	cachedMutex sync.Mutex
	cachedPath  string
)
