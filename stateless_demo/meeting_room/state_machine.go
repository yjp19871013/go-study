package meeting_room

import (
	"context"
	"fmt"
	"github.com/qmuntal/stateless"
)

const (
	triggerFree    = "free"
	triggerOccupy  = "occupy"
	triggerDisable = "disable"
)

const (
	stateFreeCode uint8 = iota + 1
	stateOccupiedCode
	stateDisabledCode
)

const (
	StateFree     = "空闲"
	StateOccupied = "占用"
	StateDisabled = "禁用"
)

var stateMap = map[uint8]string{
	stateFreeCode:     StateFree,
	stateOccupiedCode: StateOccupied,
	stateDisabledCode: StateDisabled,
}

type stateMachine struct {
	stateMachine *stateless.StateMachine
}

func newStateMachine() *stateMachine {
	sm := stateless.NewStateMachine(stateFreeCode)

	sm.Configure(stateFreeCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now meeting room is free")
			return nil
		}).
		Permit(triggerOccupy, stateOccupiedCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now meeting room will change to occupied")
			return true
		}).
		Permit(triggerDisable, stateDisabledCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now meeting room will change to disabled")
			return true
		}).
		PermitReentry(triggerFree, func(ctx context.Context, args ...interface{}) bool {
			return true
		})

	sm.Configure(stateOccupiedCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now meeting room is occupied")
			return nil
		}).
		Permit(triggerFree, stateFreeCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now meeting room will change to free")
			return true
		}).
		Permit(triggerDisable, stateDisabledCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now meeting room will change to disabled")
			return true
		})

	sm.Configure(stateDisabledCode).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("Now meeting room is disabled")
			return nil
		}).
		Permit(triggerFree, stateFreeCode, func(ctx context.Context, args ...interface{}) bool {
			fmt.Println("Now meeting room will change to free")
			return true
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

func (sm *stateMachine) Free() error {
	return sm.stateMachine.Fire(triggerFree)
}

func (sm *stateMachine) Occupied() error {
	return sm.stateMachine.Fire(triggerOccupy)
}

func (sm *stateMachine) Disable() error {
	return sm.stateMachine.Fire(triggerDisable)
}
