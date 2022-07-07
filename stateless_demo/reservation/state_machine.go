package reservation

import (
	"context"
	"fmt"
	"github.com/qmuntal/stateless"
)

const (
	triggerUsing  = "using"
	triggerUsed   = "used"
	triggerCancel = "cancel"
)

const (
	stateCreatedCode uint8 = 1 + iota
	stateUsingCode
	stateUsedCode
	stateCanceledCode
)

const (
	stateCreated  = "待使用"
	stateUsing    = "正在使用"
	stateUsed     = "已使用"
	stateCanceled = "已取消"
)

var stateMap = map[uint8]string{
	stateCreatedCode:  stateCreated,
	stateUsingCode:    stateUsing,
	stateUsedCode:     stateUsed,
	stateCanceledCode: stateCanceled,
}

type stateMachine struct {
	stateMachine *stateless.StateMachine
}

func newStateMachine() *stateMachine {
	sm := stateless.NewStateMachine(stateCreatedCode)

	sm.Configure(stateCreatedCode).
		Permit(triggerUsing, stateUsingCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now reservation will change to using")
			return true
		}).
		Permit(triggerCancel, stateCanceledCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now reservation will change to canceled")
			return true
		})

	sm.Configure(stateUsingCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now reservation is using")
			return nil
		}).
		Permit(triggerUsed, stateUsedCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now reservation will change to used")
			return true
		})

	sm.Configure(stateUsedCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now reservation is used")
			return nil
		})

	sm.Configure(stateCanceledCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now reservation is canceled")
			return nil
		})

	return &stateMachine{
		stateMachine: sm,
	}
}

func destroyStateMachine(sm *stateMachine) {
	if sm == nil {
		return
	}

	sm.stateMachine = nil
	sm = nil
}

func (sm *stateMachine) GetState() (string, error) {
	state, err := sm.stateMachine.State(context.Background())
	if err != nil {
		return "", err
	}

	return stateMap[state.(uint8)], nil
}

func (sm *stateMachine) Using() error {
	return sm.stateMachine.Fire(triggerUsing)
}

func (sm *stateMachine) Used() error {
	return sm.stateMachine.Fire(triggerUsed)
}

func (sm *stateMachine) Cancel() error {
	return sm.stateMachine.Fire(triggerCancel)
}
