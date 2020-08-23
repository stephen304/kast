package internal

import (
	"log"
	"sync"
)

type KastModule interface {
	GetName() string
	Stop() error
	Start() error
}

type DisplayMutex struct {
	m      sync.Mutex
	module KastModule
}

func NewDisplayMutex() *DisplayMutex {
	return &DisplayMutex{}
}

func (mutex *DisplayMutex) Assign(module KastModule) {
	mutex.m.Lock()
	defer mutex.m.Unlock()

	if mutex.module != module {
		if mutex.module != nil {
			log.Printf("[%s] Stopping...", mutex.module.GetName())
			go func(module KastModule) {
				module.Stop()
			}(mutex.module)
		}
		log.Printf("[%s] Starting...", module.GetName())
		mutex.module = module
		mutex.module.Start()
	}
}
