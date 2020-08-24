package internal

import (
	"log"
	"reflect"
	"sync"
)

type KastModule interface {
	Start() error
	Stop() error
}

type DisplayMutex struct {
	m      sync.Mutex
	module KastModule
}

func getModuleName(module KastModule) string {
	return reflect.Indirect(reflect.ValueOf(module)).Type().Name()
}

func NewDisplayMutex() *DisplayMutex {
	return &DisplayMutex{}
}

// Returns true if the module is different and Start was called
// On false, the module was already loaded and no action was taken
func (mutex *DisplayMutex) Assign(module KastModule) bool {
	mutex.m.Lock()
	defer mutex.m.Unlock()

	if mutex.module != module {
		if mutex.module != nil {
			log.Printf("[%s] Stopping...", getModuleName(mutex.module))
			go func(module KastModule) {
				module.Stop()
			}(mutex.module)
		}
		log.Printf("[%s] Starting...", getModuleName(module))
		mutex.module = module
		go mutex.module.Start()
		return true
	} else {
		return false
	}
}
