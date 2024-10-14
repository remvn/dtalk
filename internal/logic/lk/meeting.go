package lk

import (
	"dtalk/internal/dtalk"
	"errors"
	"sync"

	"github.com/livekit/protocol/livekit"
)

type MeetingData struct {
	m              sync.RWMutex
	roomID         string
	hostID         string
	joinRequestMap map[string]*meetingJoinRequest
}

func (data *MeetingData) GetRoomID() string {
	data.m.RLock()
	defer data.m.RUnlock()
	return data.roomID
}

func (data *MeetingData) GetHostID() string {
	data.m.RLock()
	defer data.m.RUnlock()
	return data.hostID
}

func (data *MeetingData) SetHostID(hostID string) {
	data.m.Lock()
	defer data.m.Unlock()
	data.hostID = hostID
}

func (data *MeetingData) AddJoinRequest(request *meetingJoinRequest) {
	data.m.Lock()
	defer data.m.Unlock()
	data.joinRequestMap[request.userInfo.ID] = request
}

func NewMeetingData(roomID, hostID string) *MeetingData {
	return &MeetingData{
		roomID:         roomID,
		hostID:         hostID,
		joinRequestMap: make(map[string]*meetingJoinRequest),
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
	service.meetingMap.Set(meeting.Data.GetRoomID(), meeting.Data)
	return meeting, nil
}

type meetingJoinRequest struct {
	userInfo *dtalk.UserTokenInfo
	result   chan bool
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
	if meeting.Data.GetHostID() == "" {
		return nil, ErrRoomNotReady
	}

	// handle join request
	request := &meetingJoinRequest{
		userInfo: requester,
		result:   make(chan bool),
	}
	meeting.Data.AddJoinRequest(request)

	return request.result, nil
}
