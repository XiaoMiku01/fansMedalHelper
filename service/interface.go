package service

import (
	"MedalHelper/dto"
	"sync"
)

type IConcurrency interface {
	// Exec the action of child and execute retry backup if
	Exec(user User, work sync.WaitGroup, child IExec) []dto.MedalList
}

type IExec interface {
	// Do represent real action
	Do(user User, roomID int) bool
}

// Action represent a single action for a single user
type IAction interface {
	IConcurrency
	IExec
}
