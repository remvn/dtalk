package dtalk

import "errors"

var ErrRoomNonExistent = errors.New("No longer available or non-existent room")

// this struct must be immutable
type Room struct {
	ID string
}

type JoinTokenParams struct {
	ID   string
	Name string
}
