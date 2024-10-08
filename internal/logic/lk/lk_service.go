package lk

import (
	"context"
	"dtalk/internal/pkg/cmap"
	"dtalk/internal/pkg/random"
	"errors"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Service struct {
	roomClient *lksdk.RoomServiceClient
	options    Config

	meetings *cmap.CMap[string, *Meeting]
}

type Config struct {
	HostURL   string
	ApiKey    string
	ApiSecret string
}

func NewLkService(options Config) *Service {
	roomClient := lksdk.NewRoomServiceClient(
		options.HostURL,
		options.ApiKey,
		options.ApiSecret,
	)
	service := &Service{
		roomClient: roomClient,
		options:    options,
	}

	return service
}

type JoinTokenParams struct {
	Id   string
	Name string
}

func (service *Service) GetJoinToken(roomId string, params JoinTokenParams) (string, error) {
	options := service.options
	accessToken := auth.NewAccessToken(options.ApiKey, options.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomId,
	}

	accessToken.AddGrant(grant).
		SetIdentity(params.Id).
		SetName(params.Name).
		SetValidFor(time.Hour)

	token, err := accessToken.ToJWT()
	if err != nil {
		return "", nil
	}
	return token, nil
}

func (service *Service) createRoom() (*livekit.Room, error) {
	client := service.roomClient

	roomId := random.GenerateID()
	room, err := client.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            roomId,
		EmptyTimeout:    10 * 60, // 10 minutes
		MaxParticipants: 20,
	})
	if err != nil {
		return nil, err
	}

	return room, nil
}

type Meeting struct {
	RoomId string
	HostId string
}

type MeetingOptions struct {
	HostId   string
	RoomName string
}

func (service *Service) CreateMeeting(options MeetingOptions) (*Meeting, error) {
	room, err := service.createRoom()
	if err != nil {
		return nil, err
	}

	meeting := &Meeting{
		RoomId: room.GetName(),
		HostId: options.HostId,
	}
	service.meetings.Set(meeting.RoomId, meeting)
	return meeting, nil
}

var ErrRoomNonExistent = errors.New("No longer available or non-existent room")

func (service *Service) GetRoom(roomId string) (*livekit.Room, error) {
	res, err := service.roomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{roomId},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Rooms) == 0 {
		return nil, ErrRoomNonExistent
	}
	return res.Rooms[0], nil
}

func (service *Service) ListUsers(roomId string) ([]*livekit.ParticipantInfo, error) {
	res, err := service.roomClient.ListParticipants(
		context.Background(),
		&livekit.ListParticipantsRequest{
			Room: roomId,
		},
	)
	if err != nil {
		return nil, err
	}

	return res.Participants, nil
}
