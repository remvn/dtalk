package port

import "dtalk/internal/app/dtalk"

type MeetingServiceIface interface {
	GetMeeting(roomID string) (*dtalk.Meeting, error)
	CreateMeeting(params dtalk.CreateMeetingParams) (*dtalk.Meeting, error)
	GetMeetingData(roomID string) (*dtalk.MeetingData, bool)
	AddJoinRequest(requester *dtalk.UserTokenInfo, roomID string) (<-chan bool, error)
	NotifyNewJoinRequest(roomID string) error
	GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error)
	GetParticipant(roomID string, participantID string) (*dtalk.Participant, error)
	ListParticipants(roomID string) ([]*dtalk.Participant, error)
}

type RoomClientIface interface {
	GetRoom(roomID string) (*dtalk.Room, error)
	CreateRoom() (*dtalk.Room, error)
	SendData(roomID string, destIDs []string, data any) error
	GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error)
	GetParticipant(roomID string, participantID string) (*dtalk.Participant, error)
	ListParticipants(roomID string) ([]*dtalk.Participant, error)
}
