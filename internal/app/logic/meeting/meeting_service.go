package meeting

import (
	"github.com/remvn/dtalk/internal/app/dtalk"
	"github.com/remvn/dtalk/internal/app/port"
	"github.com/remvn/dtalk/internal/pkg/cmap"
)

var _ port.MeetingServiceIface = (*MeetingService)(nil)

type MeetingService struct {
	roomClient port.RoomClientIface

	meetingMap *cmap.CMap[string, *dtalk.MeetingData]
}

func NewMeetingService(
	roomClient port.RoomClientIface,
) *MeetingService {
	return &MeetingService{
		roomClient: roomClient,
		meetingMap: cmap.New[string, *dtalk.MeetingData](),
	}
}

func (service *MeetingService) GetMeeting(roomId string) (*dtalk.Meeting, error) {
	meeting, ok := service.GetMeetingData(roomId)
	if !ok {
		return nil, dtalk.ErrRoomNonExistent
	}
	room, err := service.roomClient.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	return &dtalk.Meeting{
		Data: meeting,
		Room: room,
	}, nil
}

func (service *MeetingService) CreateMeeting(params dtalk.CreateMeetingParams) (*dtalk.Meeting, error) {
	room, err := service.roomClient.CreateRoom()
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
	service.meetingMap.Set(meeting.Data.RoomID(), meeting.Data)
	return meeting, nil
}

func (service *MeetingService) GetParticipant(roomID string, participantID string) (*dtalk.Participant, error) {
	return service.roomClient.GetParticipant(roomID, participantID)
}

func (service *MeetingService) ListParticipants(roomID string) ([]*dtalk.Participant, error) {
	return service.roomClient.ListParticipants(roomID)
}

func (service *MeetingService) GetMeetingData(roomID string) (*dtalk.MeetingData, bool) {
	return service.meetingMap.Load(roomID)
}

func (service *MeetingService) GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error) {
	return service.roomClient.GetJoinToken(roomID, params)
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

	request := &dtalk.MeetingJoinRequest{
		UserInfo:   requester,
		ResultChan: make(chan bool),
	}
	meeting.Data.AddJoinRequest(request)

	return request.ResultChan, nil
}

func (service *MeetingService) NotifyNewJoinRequest(roomID string) error {
	meeting, err := service.GetMeeting(roomID)
	if err != nil {
		return err
	}
	if meeting.Data.HostID() == "" {
		return dtalk.ErrRoomNotReady
	}

	pendingCount := len(meeting.Data.ListJoinRequesters())
	err = service.roomClient.SendData(
		roomID,
		[]string{meeting.Data.HostID()},
		NewPendingJoinRequestPacket(pendingCount),
	)
	return err
}
