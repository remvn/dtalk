package meeting

import (
	"dtalk/internal/app/dtalk"
	"dtalk/internal/app/port"
	"dtalk/internal/pkg/cmap"
)

var _ port.MeetingPort = (*MeetingService)(nil)

type MeetingService struct {
	roomManager port.RoomManager

	meetingMap *cmap.CMap[string, *dtalk.MeetingData]
}

func NewMeetingService(
	roomManager port.RoomManager,
) *MeetingService {
	return &MeetingService{
		roomManager: roomManager,
		meetingMap:  cmap.New[string, *dtalk.MeetingData](),
	}
}

func (service *MeetingService) GetMeeting(roomId string) (*dtalk.Meeting, error) {
	meeting, ok := service.Load(roomId)
	if !ok {
		return nil, dtalk.ErrRoomNonExistent
	}
	room, err := service.roomManager.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	return &dtalk.Meeting{
		Data: meeting,
		Room: room,
	}, nil
}

func (service *MeetingService) CreateMeeting(params dtalk.CreateMeetingParams) (*dtalk.Meeting, error) {
	room, err := service.roomManager.CreateRoom()
	if err != nil {
		return nil, err
	}

	meeting := &dtalk.Meeting{
		Data: dtalk.NewMeetingData(
			room.ID,
			params.RoomName,
			"",
		),
		Room: room,
	}
	service.meetingMap.Set(meeting.Data.HostID(), meeting.Data)
	return meeting, nil
}

func (service *MeetingService) GetParticipant(roomID string, participantID string) (*dtalk.Participant, error) {
	return service.roomManager.GetParticipant(roomID, participantID)
}

func (service *MeetingService) ListParticipants(roomID string) ([]*dtalk.Participant, error) {
	return service.roomManager.ListParticipants(roomID)
}

func (service *MeetingService) Load(roomID string) (*dtalk.MeetingData, bool) {
	return service.meetingMap.Load(roomID)
}

func (service *MeetingService) GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error) {
	return service.roomManager.GetJoinToken(roomID, params)
}

func (service *MeetingService) AddJoinRequest(
	requester *dtalk.UserTokenInfo,
	roomID string,
) (<-chan bool, error) {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return nil, err
	}
	if meeting.Data.HostID() == "" {
		return nil, dtalk.ErrRoomNotReady
	}

	// handle join request
	request := &dtalk.MeetingJoinRequest{
		UserInfo:   requester,
		ResultChan: make(chan bool),
	}
	meeting.Data.AddJoinRequest(request)

	return request.ResultChan, nil
}

func (service *MeetingService) SendJoinRequestPacket(roomID string) error {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return err
	}
	if meeting.Data.HostID() == "" {
		return dtalk.ErrRoomNotReady
	}

	pendingCount := len(meeting.Data.ListRequester())
	err = service.roomManager.SendData(
		roomID,
		[]string{meeting.Data.HostID()},
		NewPendingJoinRequestPacket(pendingCount),
	)
	return err
}
