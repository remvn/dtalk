package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/remvn/dtalk/internal/app/port"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoomAuth struct {
	meetingPort port.MeetingServiceIface
}

func NewRoomAuth(meetingPort port.MeetingServiceIface) *RoomAuth {
	return &RoomAuth{
		meetingPort: meetingPort,
	}
}

type RoomAuthDto struct {
	RoomID string `json:"room_id" query:"room_id"`
}

func (m *RoomAuth) Func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// can't read a Body from socket twice
		// need to re-assign a new Body after manual reading
		// TODO double check this workaround, or maybe use headers
		// instead of body
		dto := &RoomAuthDto{}
		raw, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logMiddlewareErr(c, m, err)
			return c.NoContent(http.StatusInternalServerError)
		}
		// httpserver will automatically close original
		// body, even when we re-assign with a new one (unchecked).
		// No need to close it here
		c.Request().Body = io.NopCloser(bytes.NewReader(raw))

		// log.Println("RoomAuth middleware body:", string(raw))
		jsonErr := json.Unmarshal(raw, dto)
		paramsErr := (&echo.DefaultBinder{}).BindQueryParams(c, dto)
		if jsonErr != nil && paramsErr != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		userInfo, err := ExtractUserInfo(c)
		if err != nil {
			logMiddlewareErr(c, m, err)
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

func (m *RoomAuth) Name() string {
	return "RoomAuth"
}
