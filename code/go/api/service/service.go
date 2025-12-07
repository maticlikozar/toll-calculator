package service

import (
	"sync"
)

var (
	once sync.Once

	Authorization AuthService
)

// Init func initializes used services only once.
func Init() {
	once.Do(func() {
		Authorization = Auth()
	})
}
