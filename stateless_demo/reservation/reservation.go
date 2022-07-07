package reservation

import (
	"errors"
	"stateless_demo/meeting_room"
)

type Reservation struct {
	meetingRoom *meeting_room.MeetingRoom
	*stateMachine
}

func NewReservation(meetingRoom *meeting_room.MeetingRoom) *Reservation {
	return &Reservation{
		meetingRoom:  meetingRoom,
		stateMachine: newStateMachine(),
	}
}

func DestroyReservation(reservation *Reservation) {
	if reservation == nil {
		return
	}

	reservation.meetingRoom = nil
	destroyStateMachine(reservation.stateMachine)
	reservation = nil
}

func (reservation *Reservation) Create() error {
	state, err := reservation.meetingRoom.GetState()
	if err != nil {
		return err
	}

	if state != meeting_room.StateFree {
		return errors.New("会议室不空闲，无法预约")
	}

	return nil
}

func (reservation *Reservation) Using() error {
	state, err := reservation.meetingRoom.GetState()
	if err != nil {
		return err
	}

	if state != meeting_room.StateFree {
		return errors.New("会议室不空闲，无法使用")
	}

	err = reservation.stateMachine.Using()
	if err != nil {
		return err
	}

	return reservation.meetingRoom.Occupied()
}

func (reservation *Reservation) Used() error {
	state, err := reservation.meetingRoom.GetState()
	if err != nil {
		return err
	}

	if state != meeting_room.StateOccupied {
		return errors.New("会议室未被占用")
	}

	err = reservation.stateMachine.Used()
	if err != nil {
		return err
	}

	return reservation.meetingRoom.Free()
}

func (reservation *Reservation) Cancel() error {
	state, err := reservation.meetingRoom.GetState()
	if err != nil {
		return err
	}

	if state == meeting_room.StateOccupied {
		return errors.New("会议室正被使用，无法取消")
	}

	err = reservation.stateMachine.Cancel()
	if err != nil {
		return err
	}

	return reservation.meetingRoom.Free()
}
