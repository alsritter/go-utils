package task

import (
	"fmt"
	"sync"
)

type taskServer struct {
	closeChan         chan error
	workerConcurrency int
	task              Task
	wg                sync.WaitGroup
}

func NewTaskServer(workerConcurrency int, task Task) *taskServer {
	return &taskServer{
		closeChan:         make(chan error),
		workerConcurrency: workerConcurrency,
		task:              task,
	}
}

func (l *taskServer) GracefulStop() error {
	close(l.closeChan)
	l.wg.Wait()
	l.task.Close()
	fmt.Println("graceful stop log task server")
	return nil
}

func (l *taskServer) Start() error {
	// 可以并发处理
	for i := 0; i < l.workerConcurrency; i++ {
		l.wg.Add(1)
		go l.loopHandler()
	}
	return nil
}

func (l *taskServer) loopHandler() {
	defer l.wg.Done()
	var err error
	for {
		select {
		case err = <-l.closeChan:
			fmt.Println("close log task, err:", err)
			return
		default:
		}
		l.task.Receive()
	}
}
