package main

import (
	"fmt"
	"stateless_demo/meeting_room"
	"stateless_demo/reservation"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printMeetingRoomState(meetingRoom *meeting_room.MeetingRoom) {
	state, err := meetingRoom.GetState()
	checkErr(err)

	fmt.Println("meetingRoom state:", state)
}

func printReservationState(reservation *reservation.Reservation) {
	state, err := reservation.GetState()
	checkErr(err)

	fmt.Println("reservation state:", state)
}

func main() {
	testMeetingRoom := meeting_room.NewMeetingRoom("测试")
	defer meeting_room.DestroyMeetingRoom(testMeetingRoom)

	testReservation1 := reservation.NewReservation(testMeetingRoom)
	defer reservation.DestroyReservation(testReservation1)

	err := testReservation1.Create()
	checkErr(err)
	printReservationState(testReservation1)
	printMeetingRoomState(testMeetingRoom)
	fmt.Println()

	err = testReservation1.Cancel()
	checkErr(err)
	printReservationState(testReservation1)
	printMeetingRoomState(testMeetingRoom)
	fmt.Println()

	testReservation2 := reservation.NewReservation(testMeetingRoom)
	defer reservation.DestroyReservation(testReservation2)

	err = testReservation2.Create()
	checkErr(err)
	printReservationState(testReservation2)
	printMeetingRoomState(testMeetingRoom)
	fmt.Println()

	err = testReservation2.Using()
	checkErr(err)
	printReservationState(testReservation2)
	printMeetingRoomState(testMeetingRoom)
	fmt.Println()

	err = testReservation2.Used()
	checkErr(err)
	printReservationState(testReservation2)
	printMeetingRoomState(testMeetingRoom)
	fmt.Println()
}
