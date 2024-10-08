package handler

import (
	"dtalk/internal/logic/lk"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MeetingHandler struct {
	lkService  *lk.Service
	echoServer *echo.Echo
}

func NewMeetingHandler(
	echoServer *echo.Echo,
	lkService *lk.Service,
) *MeetingHandler {
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
	Id       string `json:"id"`
	Name     string `json:"name"`
	RoomName string `json:"room_name"`
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

	meeting, err := handler.lkService.CreateMeeting(lk.MeetingOptions{
		HostId:   dto.Id,
		RoomName: dto.RoomName,
	})
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	joinToken, err := handler.lkService.GetJoinToken(meeting.RoomId, lk.JoinTokenParams{
		Id:   dto.Id,
		Name: dto.Name,
	})
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
