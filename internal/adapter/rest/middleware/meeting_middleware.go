package middleware

import (
	"dtalk/internal/app/port"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoomAuth struct {
	meetingPort port.MeetingPort
}

func NewRoomAuth(meetingPort port.MeetingPort) *RoomAuth {
	return &RoomAuth{
		meetingPort: meetingPort,
	}
}

type RoomAuthDto struct {
	RoomID string `json:"room_id" query:"room_id"`
}

func (m *RoomAuth) Func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := &RoomAuthDto{}
		if err := c.Bind(dto); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		userInfo, err := ExtractUserInfo(c)
		if err != nil {
			log.Println(fmt.Errorf("unable to process RoomAuth middleware: %w", err))
			return c.NoContent(http.StatusInternalServerError)
		}
		_, err = m.meetingPort.GetParticipant(dto.RoomID, userInfo.ID)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

func (m *RoomAuth) Apply(consumer consumer) {
	consumer.Use(m.Func)
}
