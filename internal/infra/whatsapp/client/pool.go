package client

import (
	"sync"
	"sync/atomic"
	"time"

	"go.mau.fi/whatsmeow/store"
)

// Pool gerencia pool de devices para reutilização
type Pool struct {
	devices     chan *store.Device
	maxSize     int
	maxIdle     int
	maxLifetime time.Duration
	created     int64
	reused      int64
	closed      int32
	mu          sync.RWMutex
}

// NewPool cria um novo pool de devices
func NewPool(maxSize, maxIdle int, maxLifetime time.Duration) *Pool {
	return &Pool{
		devices:     make(chan *store.Device, maxIdle),
		maxSize:     maxSize,
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
	}
}

// Get obtém um device do pool ou cria um novo
func (p *Pool) Get() *store.Device {
	if atomic.LoadInt32(&p.closed) == 1 {
		return nil
	}

	select {
	case device := <-p.devices:
		if device != nil {
			atomic.AddInt64(&p.reused, 1)
			return device
		}
	default:
		// Pool vazio, criar novo device será feito pela factory
	}

	atomic.AddInt64(&p.created, 1)
	return nil // Factory criará um novo
}

// Put retorna um device para o pool
func (p *Pool) Put(device *store.Device) {
	if device == nil || atomic.LoadInt32(&p.closed) == 1 {
		return
	}

	select {
	case p.devices <- device:
		// Device adicionado ao pool com sucesso
	default:
		// Pool cheio, descartar device
	}
}

// Close fecha o pool e limpa resources
func (p *Pool) Close() {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return // Já fechado
	}

	close(p.devices)

	// Drenar devices restantes
	for device := range p.devices {
		_ = device // Apenas drenar, devices serão coletados pelo GC
	}
}

// GetStats retorna estatísticas do pool
func (p *Pool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"maxSize":     p.maxSize,
		"maxIdle":     p.maxIdle,
		"maxLifetime": p.maxLifetime.String(),
		"available":   len(p.devices),
		"created":     atomic.LoadInt64(&p.created),
		"reused":      atomic.LoadInt64(&p.reused),
		"closed":      atomic.LoadInt32(&p.closed) == 1,
	}
}

// IsClosed verifica se o pool está fechado
func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

// Available retorna número de devices disponíveis no pool
func (p *Pool) Available() int {
	if atomic.LoadInt32(&p.closed) == 1 {
		return 0
	}
	return len(p.devices)
}

// Created retorna número total de devices criados
func (p *Pool) Created() int64 {
	return atomic.LoadInt64(&p.created)
}

// Reused retorna número total de devices reutilizados
func (p *Pool) Reused() int64 {
	return atomic.LoadInt64(&p.reused)
}
