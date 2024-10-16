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

var _ port.RoomManager = (*LkRoomManager)(nil)

type LkRoomManager struct {
	roomClient *lksdk.RoomServiceClient
	options    Config
}

type Config struct {
	HostURL   string
	ApiKey    string
	ApiSecret string
}

func NewLkRoomManager(options Config) *LkRoomManager {
	roomClient := lksdk.NewRoomServiceClient(
		options.HostURL,
		options.ApiKey,
		options.ApiSecret,
	)
	service := &LkRoomManager{
		roomClient: roomClient,
		options:    options,
	}

	return service
}

func (manager *LkRoomManager) GetJoinToken(roomID string, params dtalk.JoinTokenParams) (string, error) {
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

func (manager *LkRoomManager) CreateRoom() (*dtalk.Room, error) {
	client := manager.roomClient

	roomID := random.GenerateID()
	room, err := client.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            roomID,
		EmptyTimeout:    10 * 60, // 10 minutes
		MaxParticipants: 20,
	})

	if err != nil {
		return nil, err
	}

	return manager.convertRoom(room), nil
}

func (manager *LkRoomManager) GetRoom(roomID string) (*dtalk.Room, error) {
	res, err := manager.roomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{roomID},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Rooms) == 0 {
		return nil, dtalk.ErrRoomNonExistent
	}
	room := manager.convertRoom(res.Rooms[0])
	return room, nil
}

func (manager *LkRoomManager) convertRoom(lkRoom *livekit.Room) *dtalk.Room {
	return &dtalk.Room{
		ID: lkRoom.Name,
	}
}

func (manager *LkRoomManager) GetParticipant(roomID string, participantID string) (*dtalk.Participant, error) {
	lkPart, err := manager.getLkParticipant(roomID, participantID)
	if err != nil {
		return nil, err
	}
	return manager.convertParticipant(lkPart), nil
}

func (manager *LkRoomManager) getLkParticipant(roomID string, participantID string) (*livekit.ParticipantInfo, error) {
	res, err := manager.roomClient.GetParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     roomID,
		Identity: participantID,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (manager *LkRoomManager) ListParticipants(roomID string) ([]*dtalk.Participant, error) {
	lkArr, err := manager.listLkParticipants(roomID)
	if err != nil {
		return nil, err
	}
	arr := make([]*dtalk.Participant, 0, len(lkArr))
	for _, item := range lkArr {
		arr = append(arr, manager.convertParticipant(item))
	}
	return arr, nil
}

func (manager *LkRoomManager) listLkParticipants(roomID string) ([]*livekit.ParticipantInfo, error) {
	res, err := manager.roomClient.ListParticipants(
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

func (manager *LkRoomManager) convertParticipant(lkParticipant *livekit.ParticipantInfo) *dtalk.Participant {
	return &dtalk.Participant{
		ID:   lkParticipant.Identity,
		Name: lkParticipant.Name,
	}
}

func (manager *LkRoomManager) SendData(roomID string, destIDs []string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = manager.roomClient.SendData(context.Background(), &livekit.SendDataRequest{
		Room:                  roomID,
		Data:                  bytes,
		DestinationIdentities: destIDs,
	})
	if err != nil {
		return err
	}
	return nil
}
