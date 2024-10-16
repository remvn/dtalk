package handler

import (
	"dtalk/internal/adapter/rest/middleware"
	"dtalk/internal/app/dtalk"
	"dtalk/internal/app/port"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MeetingHandler struct {
	echoServer         *echo.Echo
	authMiddleware     *middleware.Auth
	roomAuthMiddleware *middleware.RoomAuth

	meetingPort port.MeetingServiceIface
}

func NewMeetingHandler(
	echoServer *echo.Echo,
	authMiddleware *middleware.Auth,
	roomAuthMiddlware *middleware.RoomAuth,

	meetingPort port.MeetingServiceIface,
) *MeetingHandler {
	handler := &MeetingHandler{
		echoServer:         echoServer,
		authMiddleware:     authMiddleware,
		roomAuthMiddleware: roomAuthMiddlware,

		meetingPort: meetingPort,
	}
	return handler
}

func (handler *MeetingHandler) Register(parentGroup *echo.Group) {
	const prefix = "/meeting"

	// public
	group := parentGroup.Group(prefix)
	group.POST("/create", handler.create)
	group.GET("/public-data", handler.publicData)

	// require access_token
	authGroup := parentGroup.Group(prefix)
	handler.authMiddleware.Apply(authGroup)
	authGroup.POST("/join", handler.join)

	// require access_token & belong to the room
	roomAuthGroup := parentGroup.Group(prefix)
	handler.authMiddleware.Apply(roomAuthGroup)
	handler.roomAuthMiddleware.Apply(roomAuthGroup)
	roomAuthGroup.GET("/participants", handler.listParticipants)
	roomAuthGroup.GET("/join-requesters", handler.listJoinRequesters)
	roomAuthGroup.POST("/accept", handler.accept)
}

type createMeetingDto struct {
	RoomName string `json:"room_name"`
}

type createMeetingRes struct {
	RoomID string `json:"room_id"`
}

// public routes

func (handler *MeetingHandler) create(c echo.Context) error {
	dto := &createMeetingDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	meeting, err := handler.meetingPort.CreateMeeting(dtalk.CreateMeetingParams{
		RoomName: dto.RoomName,
	})
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &createMeetingRes{
		RoomID: meeting.Data.RoomID(),
	})
}

type meetingPublicDataDto struct {
	RoomID string `query:"room_id"`
}

type meetingPublicDataRes struct {
	Name string `json:"name"`
}

func (handler *MeetingHandler) publicData(c echo.Context) error {
	dto := &meetingPublicDataDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	meeting, err := handler.meetingPort.GetMeeting(dto.RoomID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, meetingPublicDataRes{
		Name: meeting.Data.Name(),
	})
}

type joinMeetingDto struct {
	RoomID string `json:"room_id"`
}

// auth protected routes

type joinMeetingRes struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message,omitempty"`
	RoomToken string `json:"room_token,omitempty"`
}

func (handler *MeetingHandler) join(c echo.Context) error {
	dto := &joinMeetingDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userInfo, err := middleware.ExtractUserInfo(c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	meeting, err := handler.meetingPort.GetMeeting(dto.RoomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, joinMeetingRes{
			OK:      false,
			Message: dtalk.ErrRoomNonExistent.Error(),
		})
	}

	// room is just created, first one in will be the host
	if meeting.Data.HostID() == "" {
		token, err := handler.meetingPort.GetJoinToken(meeting.Data.RoomID(), dtalk.JoinTokenParams{
			ID: userInfo.ID,
		})
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		meeting.Data.SetHostID(userInfo.ID)
		return c.JSON(http.StatusOK, joinMeetingRes{
			OK:        true,
			RoomToken: token,
		})
	}

	resChan, err := handler.meetingPort.AddJoinRequest(userInfo, meeting.Data.RoomID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	_ = handler.meetingPort.NotifyNewJoinRequest(meeting.Data.RoomID())

	accepted := false
	select {
	case accepted = <-resChan:
	case <-time.After(2 * time.Minute): //  timeout
	case <-c.Request().Context().Done(): // user cancel request
	}

	log.Println("result: ", accepted)

	if accepted {
		token, err := handler.meetingPort.GetJoinToken(meeting.Data.RoomID(), dtalk.JoinTokenParams{
			ID: userInfo.ID,
		})
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, joinMeetingRes{
			OK:        true,
			RoomToken: token,
		})
	} else {
		return c.JSON(http.StatusUnauthorized, joinMeetingRes{
			OK:      false,
			Message: "Your join request gets rejected",
		})
	}
}

// auth & room protected routes

type acceptRequestDto struct {
	roomOperationDto
	Accepted    bool   `json:"accepted"`
	RequesterID string `json:"requester_id"`
}

func (handler *MeetingHandler) accept(c echo.Context) error {
	dto := &acceptRequestDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	hostInfo, err := middleware.ExtractUserInfo(c)
	if err != nil {
		logHandlerError(c, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	meeting, err := handler.meetingPort.GetMeeting(dto.RoomID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if hostInfo.ID != meeting.Data.HostID() {
		return c.NoContent(http.StatusUnauthorized)
	}

	request, ok := meeting.Data.GetJoinRequest(dto.RequesterID)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	select {
	case request.ResultChan <- dto.Accepted:
	default:
	}
	meeting.Data.RemoveJoinRequest(dto.RequesterID)

	return c.NoContent(http.StatusOK)
}

func (handler *MeetingHandler) listParticipants(c echo.Context) error {
	dto := &roomOperationDto{}
	if err := c.Bind(dto); err != nil {
		logHandlerError(c, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	arr, err := handler.meetingPort.ListParticipants(dto.RoomID)
	if err != nil {
		log.Println(fmt.Errorf("error on %s: %w", c.Path(), err))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, arr)
}

func (handler *MeetingHandler) listJoinRequesters(c echo.Context) error {
	dto := &roomOperationDto{}
	if err := c.Bind(dto); err != nil {
		logHandlerError(c, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	meeting, err := handler.meetingPort.GetMeeting(dto.RoomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, MessageRes{
			Message: dtalk.ErrRoomNonExistent.Error(),
		})
	}

	arr := meeting.Data.ListJoinRequesters()
	return c.JSON(http.StatusOK, arr)
}
