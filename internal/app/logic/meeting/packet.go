package meeting

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
	TotalCount int `json:"total_count"`
}

func NewPendingJoinRequestPacket(totalCount int) pendingJoinRequestPacket {
	return pendingJoinRequestPacket{
		packet:     packet{"new_join_request"},
		TotalCount: totalCount,
	}
}
