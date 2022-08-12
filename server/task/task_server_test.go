package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func Test_taskServer_Start(t *testing.T) {
	type fields struct {
		ts *taskServer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"简单测试",
			fields{
				ts: NewTaskServer(10,
					NewDefaultTask("testKey", redis.NewClient(&redis.Options{
						Addr:     "localhost:6379",
						Password: "123456",
						DB:       0,
					}))),
			},
			true,
		},
	}

	redis.SetLogger(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.ts.Start()

			go func() {
				for i := 0; i < 100; i++ {
					tt.fields.ts.task.Send(context.Background(), fmt.Sprintf("hello %d task server", i))
				}
			}()

			time.Sleep(2 * time.Second)

			tt.fields.ts.GracefulStop()
		})
	}
}

func Test_Channel(t *testing.T) {
	closeChan := make(chan error)

	go func(cc chan error) {
		var err error
		for {
			select {
			case err = <-cc:
				fmt.Println("close log task, err:", err)
				return
			default:
			}

			fmt.Printf("x")
		}
	}(closeChan)

	time.Sleep(time.Millisecond)

	close(closeChan)
}
