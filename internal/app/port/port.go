package port

import "dtalk/internal/app/dtalk"

type MeetingPort interface {
	GetMeeting(roomID string) (*dtalk.Meeting, error)
	CreateMeeting(params dtalk.CreateMeetingParams) (*dtalk.Meeting, error)
	Load(roomID string) (*dtalk.MeetingData, bool)
	AddJoinRequest(requester *dtalk.UserTokenInfo, roomID string) (<-chan bool, error)
	SendJoinRequestPacket(roomID string) error
	GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error)
	GetParticipant(roomID string, participantID string) (*dtalk.Participant, error)
	ListParticipants(roomID string) ([]*dtalk.Participant, error)
}

type RoomManager interface {
	GetRoom(roomID string) (*dtalk.Room, error)
	CreateRoom() (*dtalk.Room, error)
	SendData(roomID string, destIDs []string, data any) error
	GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error)
	GetParticipant(roomID string, participantID string) (*dtalk.Participant, error)
	ListParticipants(roomID string) ([]*dtalk.Participant, error)
}
