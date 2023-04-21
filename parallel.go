/**
 * @Author: Shi Jinyu
 * @Description:
 * @File:  parallel
 * @Version: 1.0.0
 * @Date: 2023/4/21 21:57
 */
package go_parallel

import (
	"context"
	"log"
	"sync"
	"time"
)

type ParallelProcessor interface {
	Do() interface{}
}

type ParallelFunc func() interface{}

func (f ParallelFunc) Do() interface{} {
	return f()
}

type processTask struct {
	name string
	p    ParallelProcessor
}

type ProcessResult struct {
	Name string
	Data interface{}
}

type ParallelObject struct {
	ctx       context.Context
	timeout   time.Duration
	cancel    func()
	wg        sync.WaitGroup
	rch       chan ProcessResult
	waitCh    chan struct{}
	endCh     chan struct{}
	p         []processTask
	r         []ProcessResult
	isTimeout bool
}

func NewParallelObject() *ParallelObject {
	return &ParallelObject{
		ctx:     context.Background(),
		timeout: 0,
		wg:      sync.WaitGroup{},
		rch:     make(chan ProcessResult),
		waitCh:  make(chan struct{}),
		endCh:   make(chan struct{}),
		p:       nil,
		r:       nil,
	}
}

func (p *ParallelObject) SetContext(ctx context.Context) *ParallelObject {
	p.ctx, p.cancel = context.WithCancel(ctx)
	return p
}

func (p *ParallelObject) SetTimeout(timeout time.Duration) *ParallelObject {
	p.timeout = timeout
	p.ctx, p.cancel = context.WithTimeout(p.ctx, timeout)
	return p
}

func (p *ParallelObject) AppendProcess(name string, f ParallelProcessor) *ParallelObject {
	p.p = append(p.p, processTask{
		name: name,
		p:    f,
	})
	return p
}

func (p *ParallelObject) AppendFunc(name string, f func() interface{}) *ParallelObject {
	p.AppendProcess(name, ParallelFunc(f))
	return p
}

func (p *ParallelObject) Run() ([]ProcessResult, bool) {
	go func() {
		defer close(p.endCh)
	loop:
		for {
			select {
			case v := <-p.rch:
				p.r = append(p.r, v)
			case <-p.ctx.Done():
				break loop
			case <-p.waitCh:
				break loop
			}
		}
		p.endCh <- struct{}{}
	}()

	go func() {
		for _, f := range p.p {
			p.wg.Add(1)
			go func(task processTask) {
				defer func() {
					if panicErr := recover(); panicErr != nil {
						log.Printf("[ERROR]:do sub process panic:%v\n", panicErr)
					}
				}()
				defer p.wg.Done()
				in := task.p.Do()
				select {
				case p.rch <- ProcessResult{
					Name: task.name,
					Data: in,
				}:
				case <-p.ctx.Done():
				}
			}(f)
		}
		p.wg.Wait()
		close(p.waitCh)
		close(p.rch)
	}()

	select {
	case <-p.waitCh:
		if p.cancel != nil {
			p.cancel()
		}
		p.isTimeout = false
	case <-p.ctx.Done():
		p.isTimeout = true
	}

	<-p.endCh
	return p.r, p.isTimeout
}
