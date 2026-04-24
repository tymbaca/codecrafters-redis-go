package service

import (
	"sync"
	"time"
)

func New() *Service {
	return &Service{
		data: make(map[string]entry),
	}
}

type Service struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	val       string
	expireSet bool
	expire    time.Time
}
