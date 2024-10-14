package lk

// type Packeter interface {
// 	GetKind() string
// }
//
// var _ Packeter = (*Packet)(nil)

type packet struct {
	Kind string `json:"kind"`
}

// func (p *Packet) GetKind() string {
// 	return p.Kind
// }

type pendingJoinRequestPacket struct {
	packet
	PendingCount int `json:"pending_count"`
}

func NewPendingJoinRequestPacket(pendingCount int) pendingJoinRequestPacket {
	return pendingJoinRequestPacket{
		packet:       packet{"new_meeting_join_request"},
		PendingCount: pendingCount,
	}
}
