package handler

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type roomOperationDto struct {
	RoomID string `json:"room_id"`
}

type MessageRes struct {
	Message string `json:"message"`
}

func logHandlerError(c echo.Context, err error) {
	log.Println(fmt.Errorf("error on %s: %w", c.Path(), err))
}
