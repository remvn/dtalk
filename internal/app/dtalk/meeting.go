package dtalk

import (
	"errors"
	"sync"
)

var ErrRoomNotReady = errors.New("This room is not ready")

type CreateMeetingParams struct {
	RoomName string
}

type Meeting struct {
	Room *Room
	Data *MeetingData
}

type MeetingData struct {
	m              sync.RWMutex
	name           string
	roomID         string
	hostID         string
	joinRequestMap map[string]*MeetingJoinRequest
}

// this struct must be immutable
type MeetingJoinRequest struct {
	UserInfo   *UserTokenInfo
	ResultChan chan bool
}

func (data *MeetingData) Name() string {
	data.m.RLock()
	defer data.m.RUnlock()
	return data.name
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

func (data *MeetingData) AddJoinRequest(request *MeetingJoinRequest) {
	data.m.Lock()
	defer data.m.Unlock()
	data.joinRequestMap[request.UserInfo.ID] = request
}

func (data *MeetingData) ListRequester() []*UserTokenInfo {
	data.m.RLock()
	defer data.m.RUnlock()
	arr := []*UserTokenInfo{}
	for _, request := range data.joinRequestMap {
		arr = append(arr, request.UserInfo)
	}
	return arr
}

func (data *MeetingData) GetJoinRequest(requesterID string) (*MeetingJoinRequest, bool) {
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

func NewMeetingData(roomID, name, hostID string) *MeetingData {
	return &MeetingData{
		roomID:         roomID,
		name:           name,
		hostID:         hostID,
		joinRequestMap: make(map[string]*MeetingJoinRequest),
	}
}
