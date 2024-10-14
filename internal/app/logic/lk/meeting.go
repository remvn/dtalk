package lk

import (
	"dtalk/internal/app/dtalk"
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

// this struct must be immutable
type meetingJoinRequest struct {
	UserInfo   *dtalk.UserTokenInfo
	ResultChan chan bool
}

func (data *MeetingData) RoomID() string {
	data.m.RLock()
	defer data.m.RUnlock()
	return data.roomID
}

func (data *MeetingData) HostID() string {
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
	data.joinRequestMap[request.UserInfo.ID] = request
}

func (data *MeetingData) ListRequester() []*dtalk.UserTokenInfo {
	data.m.RLock()
	defer data.m.RUnlock()
	arr := []*dtalk.UserTokenInfo{}
	for _, request := range data.joinRequestMap {
		arr = append(arr, request.UserInfo)
	}
	return arr
}

func (data *MeetingData) GetJoinRequest(requesterID string) (*meetingJoinRequest, bool) {
	data.m.RLock()
	defer data.m.RUnlock()
	request, ok := data.joinRequestMap[requesterID]
	return request, ok
}

func (data *MeetingData) RemoveJoinRequest(requesterID string) {
	data.m.Lock()
	defer data.m.Unlock()
	delete(data.joinRequestMap, requesterID)
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
	service.meetingMap.Set(meeting.Data.RoomID(), meeting.Data)
	return meeting, nil
}

var ErrRoomNotReady = errors.New("This room is not ready")

func (service *Service) AddJoinRequest(
	requester *dtalk.UserTokenInfo,
	roomID string,
) (<-chan bool, error) {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return nil, err
	}
	if meeting.Data.HostID() == "" {
		return nil, ErrRoomNotReady
	}

	// handle join request
	request := &meetingJoinRequest{
		UserInfo:   requester,
		ResultChan: make(chan bool),
	}
	meeting.Data.AddJoinRequest(request)

	return request.ResultChan, nil
}

func (service *Service) SendPendingJoinRequestPacket(roomID string) error {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return err
	}
	if meeting.Data.HostID() == "" {
		return ErrRoomNotReady
	}

	pendingCount := len(meeting.Data.ListRequester())
	err = service.SendData(
		roomID,
		[]string{meeting.Data.HostID()},
		NewPendingJoinRequestPacket(pendingCount),
	)
	return err
}
