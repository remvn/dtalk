package lk

import (
	"context"
	"dtalk/internal/pkg/cmap"
	"dtalk/internal/pkg/random"
	"encoding/json"
	"errors"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Service struct {
	roomClient *lksdk.RoomServiceClient
	options    Config

	meetingMap *cmap.CMap[string, *MeetingData]
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
		meetingMap: cmap.New[string, *MeetingData](),
	}

	return service
}

type JoinTokenParams struct {
	UserID   string
	UserName string
}

func (service *Service) GetJoinToken(roomId string, params JoinTokenParams) (string, error) {
	options := service.options
	accessToken := auth.NewAccessToken(options.ApiKey, options.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomId,
	}

	accessToken.AddGrant(grant).
		SetIdentity(params.UserID).
		SetName(params.UserName).
		SetValidFor(time.Hour)

	token, err := accessToken.ToJWT()
	if err != nil {
		return "", err
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

func (service *Service) SendData(roomID string, destinations []string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = service.roomClient.SendData(context.Background(), &livekit.SendDataRequest{
		Room:                  roomID,
		Data:                  bytes,
		DestinationIdentities: destinations,
	})
	if err != nil {
		return err
	}
	return nil
}
