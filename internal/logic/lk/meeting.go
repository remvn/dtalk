package lk

import (
	"dtalk/internal/dtalk"
	"errors"

	"github.com/livekit/protocol/livekit"
)

type MeetingData struct {
	RoomId string
	HostId string
}

type Meeting struct {
	Data *MeetingData
	Room *livekit.Room
}

func (service *Service) GetMeeting(roomId string) (*Meeting, error) {
	meeting, ok := service.meetingMap.Load(roomId)
	if !ok {
		return nil, ErrRoomNonExistent
	}
	room, err := service.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	return &Meeting{
		Data: meeting,
		Room: room,
	}, nil
}

type CreateMeetingParams struct {
	RoomName string
}

func (service *Service) CreateMeeting(params CreateMeetingParams) (*Meeting, error) {
	room, err := service.createRoom()
	if err != nil {
		return nil, err
	}

	meeting := &Meeting{
		Data: &MeetingData{
			RoomId: room.GetName(),
			HostId: "",
		},
		Room: room,
	}
	service.meetingMap.Set(meeting.Data.RoomId, meeting.Data)
	return meeting, nil
}

var ErrRoomNotReady = errors.New("This room is not ready")

func (service *Service) SetupMeetingJoinRequest(
	requester *dtalk.UserTokenInfo,
	roomID string,
) (chan<- string, error) {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return nil, err
	}
	if meeting.Data.HostId == "" {
		return nil, ErrRoomNotReady
	}
	return nil, nil
}
