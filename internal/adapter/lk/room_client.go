package lk

import (
	"context"
	"dtalk/internal/app/dtalk"
	"dtalk/internal/app/port"
	"dtalk/internal/pkg/random"
	"encoding/json"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

var _ port.RoomClientIface = (*RoomClient)(nil)

type RoomClient struct {
	service *lksdk.RoomServiceClient
	options Config
}

type Config struct {
	HostURL   string
	ApiKey    string
	ApiSecret string
}

func NewRoomClient(options Config) *RoomClient {
	service := lksdk.NewRoomServiceClient(
		options.HostURL,
		options.ApiKey,
		options.ApiSecret,
	)
	client := &RoomClient{
		service: service,
		options: options,
	}

	return client
}

func (manager *RoomClient) GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error) {
	options := manager.options
	token := auth.NewAccessToken(options.ApiKey, options.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}

	token.AddGrant(grant).
		SetIdentity(params.ID).
		SetName(params.Name).
		SetValidFor(time.Hour)

	tokenStr, err := token.ToJWT()
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (client *RoomClient) CreateRoom() (*dtalk.Room, error) {
	service := client.service

	roomID := random.GenerateID()
	room, err := service.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            roomID,
		EmptyTimeout:    10 * 60, // 10 minutes
		MaxParticipants: 20,
	})

	if err != nil {
		return nil, err
	}

	return client.convertRoom(room), nil
}

func (client *RoomClient) GetRoom(roomID string) (*dtalk.Room, error) {
	res, err := client.service.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{roomID},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Rooms) == 0 {
		return nil, dtalk.ErrRoomNonExistent
	}
	room := client.convertRoom(res.Rooms[0])
	return room, nil
}

func (client *RoomClient) convertRoom(lkRoom *livekit.Room) *dtalk.Room {
	return &dtalk.Room{
		ID: lkRoom.Name,
	}
}

func (client *RoomClient) GetParticipant(roomID string, participantID string) (*dtalk.Participant, error) {
	lkPart, err := client.getLkParticipant(roomID, participantID)
	if err != nil {
		return nil, err
	}
	return client.convertParticipant(lkPart), nil
}

func (client *RoomClient) getLkParticipant(roomID string, participantID string) (*livekit.ParticipantInfo, error) {
	res, err := client.service.GetParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     roomID,
		Identity: participantID,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *RoomClient) ListParticipants(roomID string) ([]*dtalk.Participant, error) {
	lkArr, err := client.listLkParticipants(roomID)
	if err != nil {
		return nil, err
	}
	arr := make([]*dtalk.Participant, 0, len(lkArr))
	for _, item := range lkArr {
		arr = append(arr, client.convertParticipant(item))
	}
	return arr, nil
}

func (client *RoomClient) listLkParticipants(roomID string) ([]*livekit.ParticipantInfo, error) {
	res, err := client.service.ListParticipants(
		context.Background(),
		&livekit.ListParticipantsRequest{
			Room: roomID,
		},
	)
	if err != nil {
		return nil, err
	}
	return res.Participants, nil
}

func (client *RoomClient) convertParticipant(lkParticipant *livekit.ParticipantInfo) *dtalk.Participant {
	return &dtalk.Participant{
		ID:   lkParticipant.Identity,
		Name: lkParticipant.Name,
	}
}

func (client *RoomClient) SendData(roomID string, destIDs []string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = client.service.SendData(context.Background(), &livekit.SendDataRequest{
		Room:                  roomID,
		Data:                  bytes,
		DestinationIdentities: destIDs,
	})
	if err != nil {
		return err
	}
	return nil
}
