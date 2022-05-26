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

// Action represent single action for single user
type IAction interface{
	// Exec the action, if fail execute retry backup
	Exec(user User, medals []dto.MedalList, work sync.WaitGroup) []dto.MedalList
	// Do represent real action
	Do(user User, roomID int) bool
}

type Action struct {}

func (a Action) Exec(user User, medals []dto.MedalList, work sync.WaitGroup) []dto.MedalList {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	fail := make([]dto.MedalList, 0, len(medals))
	for _, medal := range medals {
		wg.Add(1)
		backup := retry.NewFibonacci(1 * time.Second)
		backup = retry.WithMaxRetries(uint64(user.retryTimes), backup)
		go func(medal dto.MedalList) {
			ctx := context.Background()
			err := retry.Do(ctx, backup, func(ctx context.Context) error {
				if ok := a.Do(user, medal.RoomInfo.RoomID); !ok {
					return retry.RetryableError(errors.New("Action fail"))
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

// Do is specific thing for each action
func (Action) Do(user User, roomID int) bool {
	return true
}

// Like include 3 * like
type Like struct {
	Action
}

func (Like) Do(user User, roomID int) bool {
	return manager.LikeInteract(user.accessKey, roomID)
}

// Share include 5 * share
type Share struct {
	Action
}

func (Share) Do(user User, roomID int) bool {
	return manager.ShareRoom(user.accessKey, roomID)
}

// Danmaku include sending daily danmu
type Danmaku struct {
	Action
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
		go action.Exec(task.User, task.medalsLow, wg)
	}
	wg.Wait()
}
