package service

import (
	"MedalHelper/dto"
	"MedalHelper/manager"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sethvargo/go-retry"
)

// SyncAction implement IConcurrency, support synchronous actions
type SyncAction struct{}

func (a *SyncAction) Exec(user User, work sync.WaitGroup, child IExec) []dto.MedalList {
	fail := make([]dto.MedalList, 0, len(user.medalsLow))
	for _, medal := range user.medalsLow {
		backup := retry.NewFibonacci(1 * time.Second)
		backup = retry.WithMaxRetries(uint64(user.retryTimes), backup)
		ctx := context.Background()
		err := retry.Do(ctx, backup, func(ctx context.Context) error {
			if ok := child.Do(user, medal.RoomInfo.RoomID); !ok {
				return retry.RetryableError(errors.New("action fail"))
			}
			return nil
		})
		if err != nil {
			fail = append(fail, medal)
		}
	}
	work.Done()
	return fail
}

// AsyncAction implement IConcurrency, support asynchronous actions
type AsyncAction struct{}

func (a *AsyncAction) Exec(user User, work sync.WaitGroup, child IExec) []dto.MedalList {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	fail := make([]dto.MedalList, 0, len(user.medalsLow))
	for _, medal := range user.medalsLow {
		wg.Add(1)
		backup := retry.NewFibonacci(1 * time.Second)
		backup = retry.WithMaxRetries(uint64(user.retryTimes), backup)
		go func(medal dto.MedalList) {
			ctx := context.Background()
			err := retry.Do(ctx, backup, func(ctx context.Context) error {
				if ok := child.Do(user, medal.RoomInfo.RoomID); !ok {
					return retry.RetryableError(errors.New("action fail"))
				}
				return nil
			})
			if err != nil {
				mu.Lock()
				fail = append(fail, medal)
				mu.Unlock()
			}
			wg.Done()
		}(medal)
	}
	wg.Wait()
	work.Done()
	return fail
}

// Like implement IExec, include 3 * like
type Like struct {
	AsyncAction
}

func (Like) Do(user User, roomID int) bool {
	times := 3
	for i := 0; i < times; i++ {
		if ok := manager.LikeInteract(user.accessKey, roomID); !ok {
			return false
		}
	}
	user.info("点赞完成 %d", roomID)
	return true
}

// Share implement IExec, include 5 * share
type Share struct {
	SyncAction
}

func (Share) Do(user User, roomID int) bool {
	times := 5
	for i := 0; i < times; i++ {
		if ok := manager.ShareRoom(user.accessKey, roomID); !ok {
			return false
		}
		// FIXME: how long is waiting time for share?
		time.NewTimer(1 * time.Second)
	}
	user.info("分享完成 %d", roomID)
	return true
}

// Danmaku implement IExec, include sending daily danmu
type Danmaku struct {
	SyncAction
}

func (Danmaku) Do(user User, roomID int) bool {
	return manager.SendDanmaku(user.accessKey, roomID)
}

// Task aggregate user info and corresponding action
type Task struct {
	User
	actions []IAction
}

func NewTask(user User, actions []IAction) Task {
	return Task{
		User:    user,
		actions: actions,
	}
}

func (task *Task) Start() {
	wg := sync.WaitGroup{}
	for _, action := range task.actions {
		wg.Add(1)
		go action.Exec(task.User, wg, action)
	}
	wg.Wait()
}
