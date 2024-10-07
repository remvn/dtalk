package handler

import (
	"dtalk/internal/dtalk"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MeetingHandler struct {
	lkService  *dtalk.LkService
	echoServer *echo.Echo
}

func NewMeetingHandler(echoServer *echo.Echo, lkService *dtalk.LkService) *MeetingHandler {
	handler := &MeetingHandler{
		echoServer: echoServer,
		lkService:  lkService,
	}
	return handler
}

func (handler *MeetingHandler) Register() {
	group := handler.echoServer.Group("meeting")
	group.POST("/create", handler.create)
	group.POST("/join-request", handler.joinRequest)
}

type createMeetingDto struct {
	RoomName string `json:"room_name"`
	HostId   string `json:"host_id"`
}

type createMeetingRes struct {
	JoinToken string `json:"join_token"`
}

func (handler *MeetingHandler) create(c echo.Context) error {
	dto := &createMeetingDto{}
	if err := c.Bind(dto); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	meeting, err := handler.lkService.CreateMeeting(dtalk.MeetingOptions{
		HostId:   dto.HostId,
		RoomName: dto.RoomName,
	})
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	joinToken, err := handler.lkService.GetJoinToken(meeting.RoomId, dto.HostId)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &createMeetingRes{
		JoinToken: joinToken,
	})
}

func (handler *MeetingHandler) joinRequest(c echo.Context) error {
	return nil
}
