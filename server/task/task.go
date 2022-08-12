package task

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

// 任务接口
type Task interface {
	// 接受处理任务
	Receive()

	// 发送任务
	Send(ctx context.Context, data interface{})

	Close() error
}

type defaultTask struct {
	rds *redis.Client
	key string
}

func NewDefaultTask(key string, rds *redis.Client) Task {
	return &defaultTask{key: key, rds: rds}
}

func (t *defaultTask) Receive() {
	data := t.rds.BRPop(1, t.key)
	if err := data.Err(); err != nil {
		panic(err)
	}

	fmt.Println(data.Val()[1])
}

func (t *defaultTask) Send(ctx context.Context, data interface{}) {
	if err := t.rds.LPush(t.key, data).Err(); err != nil {
		panic(err)
	}
}

func (t *defaultTask) Close() error {
	return t.rds.Close()
}
