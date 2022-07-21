package core

import (
	"osdsvr/pkg/zlog"
	"sync"

	"go.uber.org/zap"
)

type IDGenerator interface {
	GenerateOid(obj ...interface{}) (int64, error)
}

type IncrementIDGenerator struct {
	mu       sync.Mutex
	activeID int64
}

func (g *IncrementIDGenerator) Serve() error {
	if err := g.LoadActiveID(); err != nil {
		zlog.Error("load active id failed")
		return err
	}
	return nil
}

func (g *IncrementIDGenerator) GenerateOid(obj ...interface{}) (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.activeID++
	// todo： store id
	return g.activeID, nil
}

func (g *IncrementIDGenerator) LoadActiveID() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	zlog.Info("Increment IDGenerator load ActiveID successfully", zap.Int64("activeID", g.activeID))
	return nil
}
