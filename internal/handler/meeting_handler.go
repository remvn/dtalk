package handler

import (
	"dtalk/internal/logic/lk"
	"dtalk/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MeetingHandler struct {
	lkService      *lk.Service
	echoServer     *echo.Echo
	authMiddleware *middleware.AuthMiddleware
}

func NewMeetingHandler(
	echoServer *echo.Echo,
	lkService *lk.Service,
	authMiddleware *middleware.AuthMiddleware,
) *MeetingHandler {
	handler := &MeetingHandler{
		echoServer:     echoServer,
		lkService:      lkService,
		authMiddleware: authMiddleware,
	}
	return handler
}

func (handler *MeetingHandler) Register() {
	protectedGroup := handler.echoServer.Group("meeting")
	handler.authMiddleware.Apply(protectedGroup)
	protectedGroup.POST("/join", handler.join)
	protectedGroup.POST("/accept", handler.accept)

	group := handler.echoServer.Group("meeting")
	group.POST("/create", handler.create)
}

type createMeetingDto struct {
	RoomName string `json:"room_name"`
}

type createMeetingRes struct {
	RoomID string `json:"room_id"`
}

func (handler *MeetingHandler) create(c echo.Context) error {
	dto := &createMeetingDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	meeting, err := handler.lkService.CreateMeeting(lk.CreateMeetingParams{
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

type joinMeetingDto struct {
	RoomID string `json:"room_id"`
}

type joinMeetingRes struct {
	OK          bool   `json:"ok"`
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
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

	meeting, err := handler.lkService.GetMeeting(dto.RoomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, joinMeetingRes{
			OK:      false,
			Message: lk.ErrRoomNonExistent.Error(),
		})
	}

	// room is just created, first one in will be the host
	if meeting.Data.HostID() == "" {
		token, err := handler.lkService.GetJoinToken(meeting.Data.RoomID(), lk.JoinTokenParams{
			UserID: userInfo.ID,
		})
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		meeting.Data.SetHostID(userInfo.ID)
		return c.JSON(http.StatusOK, joinMeetingRes{
			OK:          true,
			AccessToken: token,
		})
	}

	resChan, err := handler.lkService.AddJoinRequest(userInfo, meeting.Data.RoomID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	_ = handler.lkService.SendPendingJoinRequestPacket(meeting.Data.RoomID())

	accepted := false
	select {
	case accepted = <-resChan:
	case <-time.After(2 * time.Minute): //  timeout
	case <-c.Request().Context().Done(): // user cancel request
	}

	log.Println("result: ", accepted)

	if accepted {
		token, err := handler.lkService.GetJoinToken(meeting.Data.RoomID(), lk.JoinTokenParams{
			UserID: userInfo.ID,
		})
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, joinMeetingRes{
			OK:          true,
			AccessToken: token,
		})
	} else {
		return c.JSON(http.StatusUnauthorized, joinMeetingRes{
			OK:      false,
			Message: "Your join request gets rejected",
		})
	}
}

type acceptRequestDto struct {
	Accepted    bool   `json:"accepted"`
	RoomID      string `json:"room_id"`
	RequesterID string `json:"requester_id"`
}

func (handler *MeetingHandler) accept(c echo.Context) error {
	dto := &acceptRequestDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	hostInfo, err := middleware.ExtractUserInfo(c)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	meeting, err := handler.lkService.GetMeeting(dto.RoomID)
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
