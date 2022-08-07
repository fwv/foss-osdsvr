package core

import (
	"errors"
	"osdsvr/pkg/zlog"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var (
	CLOSED  = 1
	RUNNING = 0
)

type Task struct {
	DoTask func(v ...interface{}) error
	Param  []interface{}
}

type Scheduler struct {
	mu               sync.Mutex
	Capacity         int64
	taskCh           chan *Task
	status           int
	runningWorkerNum int64
	taskNum          int64
}

func NewScheduler(capacity int64) *Scheduler {
	return &Scheduler{
		Capacity:         capacity,
		taskCh:           make(chan *Task, capacity),
		status:           RUNNING,
		runningWorkerNum: 0,
	}
}

func (s *Scheduler) AddTask(t *Task) error {
	if s.getStatus() == CLOSED {
		return errors.New("scheduler is closed")
	}
	s.taskCh <- t
	s.incTaskNum()
	return nil
}

func (s *Scheduler) Start(concurrency int64) error {
	zlog.Info("scheduler start working", zap.Any("concurrency", concurrency))
	for i := 0; i < int(concurrency); i++ {
		if err := s.LaunchWorker(); err != nil {
			zlog.Error("launch worker failed", zap.Error(err))
		}
	}
	return nil
}

// start a worker to handle task
func (s *Scheduler) LaunchWorker() error {
	s.incRunningWorkerNum()
	go func() {
		defer s.decRunningWorkerNum()
		for {
			zlog.Debug("task queue", zap.Any("size", s.getTaskNum()))
			select {
			case task, ok := <-s.taskCh:
				if !ok {
					break
				}
				if err := task.DoTask(task.Param...); err != nil {
					zlog.Error("failed to do task", zap.Error(err))
					break
				}
				s.decTaskNum()
			}
		}
	}()
	return nil
}

func (s *Scheduler) Shutdown() error {
	s.setStatus(CLOSED)
	for s.getTaskNum() > 0 {
		time.Sleep(time.Second)
	}
	close(s.taskCh)
	return nil
}

func (s *Scheduler) incRunningWorkerNum() {
	atomic.AddInt64(&s.runningWorkerNum, 1)
	// zlog.Info("num of running worker inc", zap.Int64("num", newV))
}

func (s *Scheduler) decRunningWorkerNum() {
	atomic.AddInt64(&s.runningWorkerNum, -1)
	// zlog.Info("num of running worker decrease", zap.Int64("num", newV))
}

func (s *Scheduler) getRunningWorkerNum() int64 {
	return atomic.LoadInt64(&s.runningWorkerNum)
}

func (s *Scheduler) incTaskNum() {
	atomic.AddInt64(&s.taskNum, 1)
}

func (s *Scheduler) decTaskNum() {
	atomic.AddInt64(&s.taskNum, -1)
}

func (s *Scheduler) getTaskNum() int64 {
	val := atomic.LoadInt64(&s.taskNum)
	return val
}

func (s *Scheduler) getStatus() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *Scheduler) setStatus(status int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
}
