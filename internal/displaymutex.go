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
			go func(module KastModule) {
				log.Printf("Stopping module: %s", module.GetName())
				module.Stop()
			}(mutex.module)
		}
		log.Printf("Loading module: %s", module.GetName())
		mutex.module = module
		log.Printf("Starting module: %s", module.GetName())
		mutex.module.Start()
	}
}
