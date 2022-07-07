package meeting_room

type MeetingRoom struct {
	name string
	*stateMachine
}

func NewMeetingRoom(name string) *MeetingRoom {
	return &MeetingRoom{
		name:         name,
		stateMachine: newStateMachine(),
	}
}

func DestroyMeetingRoom(meetingRoom *MeetingRoom) {
	if meetingRoom == nil {
		return
	}

	meetingRoom.name = ""
	destroyStateMachine(meetingRoom.stateMachine)
	meetingRoom = nil
}

func (meetingRoom *MeetingRoom) GetState() (string, error) {
	return meetingRoom.stateMachine.GetState()
}
