package lk

import (
	"dtalk/internal/dtalk"
	"dtalk/internal/pkg/cmap"
	"errors"

	"github.com/livekit/protocol/livekit"
)

type MeetingData struct {
	RoomID         string
	HostID         string
	JoinRequestMap *cmap.CMap[string, *meetingJoinRequest]
}

func NewMeetingData(roomID, hostID string) *MeetingData {
	return &MeetingData{
		RoomID:         roomID,
		HostID:         hostID,
		JoinRequestMap: cmap.New[string, *meetingJoinRequest](),
	}
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
		Data: NewMeetingData(room.Name, ""),
		Room: room,
	}
	service.meetingMap.Set(meeting.Data.RoomID, meeting.Data)
	return meeting, nil
}

type meetingJoinRequest struct {
	userInfo *dtalk.UserTokenInfo
	result   chan<- bool
}

var ErrRoomNotReady = errors.New("This room is not ready")

func (service *Service) SetupMeetingJoinRequest(
	requester *dtalk.UserTokenInfo,
	roomID string,
) (<-chan bool, error) {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return nil, err
	}
	if meeting.Data.HostID == "" {
		return nil, ErrRoomNotReady
	}
	request := &meetingJoinRequest{
		userInfo: requester,
		result:   make(chan bool),
	}
	meeting.Data.JoinRequestMap.Set(requester.ID, request)
	return nil, nil
}
